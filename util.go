package mapreduce

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/BWbwchen/MapReduce/master"
	"github.com/BWbwchen/MapReduce/worker"
	"github.com/spf13/cobra"
)

func ParseArg() ([]string, string, int, int, bool) {
	var files []string
	var nReducer int64
	var nWorker int64
	var plugin string
	var inRAM bool
	var port int64
	var rootCmd = &cobra.Command{
		Use:   "mapreduce",
		Short: "MapReduce is an easy-to-use parallel framework by Bo-Wei Chen(BWbwchen)",
		Long: `MapReudce is an easy-to-use Map Reduce Go parallel-computing framework inspired by 2021 6.824 lab1.
It supports multiple workers threads on a single machine and multiple processes on a single machine right now.`,
		Run: func(cmd *cobra.Command, args []string) {
			tempFiles := []string{}

			for _, f := range files {
				// expand the file path
				expandFiles, err := filepath.Glob(f)
				if err != nil {
					panic(err)
				}
				tempFiles = append(tempFiles, expandFiles...)
			}

			pluginFiles, err := filepath.Glob(plugin)
			if err != nil {
				panic(err)
			} else if len(pluginFiles) == 0 {
				panic("No such file")
			}

			plugin = pluginFiles[0]
			files = tempFiles
			MasterIP = ":" + strconv.Itoa(int(port))
		},
	}

	rootCmd.PersistentFlags().StringSliceVarP(&files, "input", "i", []string{}, "Input files")
	rootCmd.MarkPersistentFlagRequired("input")
	rootCmd.PersistentFlags().StringVarP(&plugin, "plugin", "p", "", "Plugin .so file")
	rootCmd.MarkPersistentFlagRequired("plugin")
	rootCmd.PersistentFlags().Int64VarP(&nReducer, "reduce", "r", 1, "Number of Reducers")
	rootCmd.PersistentFlags().Int64VarP(&nWorker, "worker", "w", 4, "Number of Workers(for master node)\nID of worker(for worker node)")
	rootCmd.PersistentFlags().Int64Var(&port, "port", 10000, "Port number")
	rootCmd.PersistentFlags().BoolVarP(&inRAM, "inRAM", "m", true, "Whether write the intermediate file in RAM")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return files, plugin, int(nReducer), int(nWorker), inRAM
}

func startSingleMachineWorker(plugin string, nWorker int, nReducer int, storeInRAM bool) {
	if nWorker < nReducer {
		panic("Need more worker!")
	}

	pluginFile, _ := filepath.Abs(plugin)

	var wg sync.WaitGroup
	worker.Init(MasterIP)

	// Start Worker
	for i := 0; i < nWorker; i++ {
		wg.Add(1)
		go func(i0 int) {
			worker.StartWorker(pluginFile, nReducer, fmt.Sprintf(":1000%v", i0+1), storeInRAM)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func startMaster(input []string, nWorker int, nReducer int) {
	inputFiles := []string{}
	for _, s := range input {
		f, _ := filepath.Abs(s)
		inputFiles = append(inputFiles, f)
	}

	var wg sync.WaitGroup
	// Start master
	// master.StartMaster(os.Args[1:], nReducer, MasterIP)
	wg.Add(1)
	go func() {
		master.StartMaster(inputFiles, nWorker, nReducer, MasterIP)
		wg.Done()
	}()

	wg.Wait()
}

func startWorker(plugin string, id int, nReducer int, storeInRAM bool) {
	pluginFile, _ := filepath.Abs(plugin)

	var wg sync.WaitGroup
	worker.Init(MasterIP)

	// Start Worker
	wg.Add(1)
	go func() {
		worker.StartWorker(pluginFile, nReducer, fmt.Sprintf(":1000%v", id+1), storeInRAM)
		wg.Done()
	}()

	wg.Wait()
}
