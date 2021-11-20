package main

import (
	"fmt"
	"path/filepath"
	"sync"
	// "time"

	// "os"

	"github.com/BWbwchen/MapReduce/master"
	"github.com/BWbwchen/MapReduce/worker"
)

var MasterIP string = ":10000"

func main() {
	var input []string
	nReducer := 2
	nWorker := 2

	input = []string{"txt/pg-being_ernest.txt", "txt/pg-grimm.txt"}
	job(input, nWorker, nReducer, "cmd/wc.so")

	input = []string{}
	for i := 0; i < nReducer; i++ {
		input = append(input, fmt.Sprintf("output/mr-out-%v.txt", i))
	}
	nReducer = 1
	nWorker = 2
	job(input, nWorker, nReducer, "cmd/merge.so")
}

func job(input []string, nWorker int, nReducer int, plugin string) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		startMaster(input, nReducer)
		wg.Done()
	}()

	// time.Sleep(2 * time.Second)
	wg.Add(1)
	go func() {
		startWorker(plugin, nWorker, nReducer)
		wg.Done()
	}()

	wg.Wait()
}

func startWorker(plugin string, nWorker int, nReducer int) {
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
			worker.StartWorker(pluginFile, nReducer, fmt.Sprintf(":1000%v", i0+1))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func startMaster(input []string, nReducer int) {
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
		master.StartMaster(inputFiles, nReducer, MasterIP)
		wg.Done()
	}()

	wg.Wait()
}
