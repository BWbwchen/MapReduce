package mapreduce

func StartWorker(input []string, plugin string, nReducer int, nWorker int, storeInRAM bool) {
	startWorker(plugin, nWorker, nReducer, storeInRAM)
}
