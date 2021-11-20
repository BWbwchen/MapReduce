package mapreduce

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	// "time"

	// "os"

	"github.com/BWbwchen/MapReduce/master"
	"github.com/BWbwchen/MapReduce/worker"
	"github.com/spf13/cobra"
)

var MasterIP string = ":10000"

func StartSingleMachineJob(input []string, plugin string, nReducer int, nWorker int, inRAM bool) {
	if len(input) == 0 {
		os.Exit(0)
	}
	job(input, nWorker, nReducer, plugin, inRAM)
}

func ParseArg() ([]string, string, int, int, bool) {
	var files []string
	var nReducer int64
	var nWorker int64
	var plugin string
	var inRAM bool
	var port int64
	var rootCmd = &cobra.Command{
		Use:   "mapreduce",
		Short: "mapreduce is a easy-to-use parallel framework",
		Long:  `mapreduce is a easy-to-use parallel framework`,
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
	rootCmd.PersistentFlags().Int64VarP(&nWorker, "worker", "w", 4, "Number of Workers")
	rootCmd.PersistentFlags().Int64Var(&port, "port", 10000, "Port number")
	rootCmd.PersistentFlags().BoolVarP(&inRAM, "inRAM", "m", true, "Whether write the intermediate file in RAM")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return files, plugin, int(nReducer), int(nWorker), inRAM
}

func job(input []string, nWorker int, nReducer int, plugin string, storeInRAM bool) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		startMaster(input, nWorker, nReducer)
		wg.Done()
	}()

	// time.Sleep(2 * time.Second)
	wg.Add(1)
	go func() {
		startWorker(plugin, nWorker, nReducer, storeInRAM)
		wg.Done()
	}()

	wg.Wait()
}

func startWorker(plugin string, nWorker int, nReducer int, storeInRAM bool) {
	if nWorker < nReducer {
		panic("Need more worker!")
	}
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
	// 	os.Exit(1)
	// }
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
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
	// 	os.Exit(1)
	// }

	inputFiles := []string{}
	for _, s := range input {
		f, _ := filepath.Abs(s)
		inputFiles = append(inputFiles, f)
	}
	// pluginFile := "../wc.so"

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
