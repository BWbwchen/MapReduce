package main

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	// "os"

	"github.com/BWbwchen/MapReduce/master"
	"github.com/BWbwchen/MapReduce/worker"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		startMaster()
		wg.Done()
	}()

	time.Sleep(2 * time.Second)
	wg.Add(1)
	go func() {
		startWorker()
		wg.Done()
	}()

	wg.Wait()
}

func startWorker() {
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
	// 	os.Exit(1)
	// }
	pluginFile, _ := filepath.Abs("cmd/wc.so")

	nReducer := 1
	MasterIP := ":10000"

	var wg sync.WaitGroup
	worker.Init(MasterIP)

	// Start Worker
	for i := 0; i < nReducer; i++ {
		wg.Add(1)
		go func(i0 int) {
			worker.StartWorker(pluginFile, nReducer, fmt.Sprintf(":1000%v", i0+1))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func startMaster() {
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
	// 	os.Exit(1)
	// }

	input := []string{ /* "txt/pg-being_ernest.txt",  */ "txt/pg-grimm.txt"}
	inputFiles := []string{}
	for _, s := range input {
		f, _ := filepath.Abs(s)
		inputFiles = append(inputFiles, f)
	}
	// pluginFile := "../wc.so"

	nReducer := 1
	MasterIP := ":10000"

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
