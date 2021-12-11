package mapreduce

import (
	"os"
	"sync"
)

var MasterIP string = ":10000"

func StartSingleMachineJob(input []string, plugin string, nReducer int, nWorker int, inRAM bool) {
	if len(input) == 0 {
		os.Exit(0)
	}
	singleMachineJob(input, nWorker, nReducer, plugin, inRAM)
}

func singleMachineJob(input []string, nWorker int, nReducer int, plugin string, storeInRAM bool) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		startMaster(input, nWorker, nReducer)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		startSingleMachineWorker(plugin, nWorker, nReducer, storeInRAM)
		wg.Done()
	}()

	wg.Wait()
}
