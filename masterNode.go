package mapreduce

func StartMaster(input []string, plugin string, nReducer int, nWorker int, inRAM bool) {
	startMaster(input, nWorker, nReducer)
}
