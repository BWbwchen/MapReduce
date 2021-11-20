package master

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
)

type Master struct {
	Workers      []WorkerInfo
	MapTasks     []MapTaskInfo
	ReduceTasks  []ReduceTaskInfo
	numWorkers   int
	totalWorkers int
	numReducer   int
	enoughWorker chan bool
	mux          sync.Mutex
	rpc.UnimplementedMasterServer
}

func NewMaster(nWorker int, nReduce int) rpc.MasterServer {
	return &Master{
		ReduceTasks:  newReduceTasks(nReduce),
		numWorkers:   0,
		totalWorkers: nWorker,
		numReducer:   nReduce,
		enoughWorker: make(chan bool, 1),
	}
}

// gRPC functions

func (ms *Master) WorkerRegister(ctx context.Context, in *rpc.WorkerInfo) (*rpc.RegisterResult, error) {
	var num int
	nw := NewWorker(in.Uuid, in.Ip)
	ms.mux.Lock()
	ms.Workers = append(ms.Workers, nw)
	ms.numWorkers++

	num = ms.numWorkers
	ms.mux.Unlock()
	log.Info("[Master] Worker register success")
	return &rpc.RegisterResult{Result: true, Id: int64(num - 1)}, nil
}

func (ms *Master) UpdateIMDInfo(ctx context.Context, in *rpc.IMDInfo) (*rpc.UpdateResult, error) {
	for i, f := range in.Filenames {
		ms.mux.Lock()
		ms.ReduceTasks[i].IMDs = append(ms.ReduceTasks[i].IMDs,
			IMDInfo{
				IP:       ms.ServiceDiscovey(in.Uuid),
				FileName: f,
			})
		ms.mux.Unlock()
	}
	log.Info(fmt.Sprintf("[Master] %v update IMD info success", in.Uuid))
	return &rpc.UpdateResult{Result: true}, nil
}

func (ms *Master) ServiceDiscovey(uuid string) string {
	var ip string

	for _, wi := range ms.Workers {
		if wi.UUID == uuid {
			ip = wi.GetIP()
		}
	}

	return ip
}

// Normal Functions

func (ms *Master) WaitForEnoughWorker() {
	log.Trace("[Master] Wait for enough workers")
	nWorker, totalWorker := ms.getWorkerNum()
	for nWorker < totalWorker {
		nWorker, totalWorker = ms.getWorkerNum()
	}
	log.Trace("[Master] Enough workers!")
}

func (ms *Master) getWorkerNum() (int, int) {
	ms.mux.Lock()
	nWorkers := ms.numWorkers
	tWorkers := ms.totalWorkers
	ms.mux.Unlock()
	return nWorkers, tWorkers
}

func (ms *Master) DistributeWork(files []string) {
	log.Trace("[Master] Start distribute workload")
	_, numWorkers := ms.AvailableWorkers()
	// Initialize MapTasks
	ms.MapTasks = newMapTasks(numWorkers)

	// Distribute work
	// Count the total lines
	for _, file := range files {
		totalLine := lineNums(file)
		baseWorkLoad := totalLine / numWorkers

		from := 0
		for i := 0; i < numWorkers; i++ {
			workLoad := baseWorkLoad
			if i < (totalLine % numWorkers) {
				workLoad++
			}

			ms.MapTasks[i].AddFile(file, from, from+workLoad)
			from += workLoad
		}
	}
	log.Trace("[Master] End distribute workload")
}

func (ms *Master) AvailableWorkers() ([]WorkerInfo, int) {
	ms.mux.Lock()
	retInfo := ms.Workers
	retWorkers := ms.numWorkers
	ms.mux.Unlock()
	return retInfo, retWorkers
}

func lineNums(file string) int {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	num := 0
	for scanner.Scan() {
		num++
	}
	return num
}

func (ms *Master) DistributeMapTask() {
	log.Trace("[Master] Start Map task")
	var wg sync.WaitGroup
	workerID := 0
	workers, nWorkers := ms.AvailableWorkers()
	for _, mapTask := range ms.MapTasks {
		wg.Add(1)
		go func(task MapTaskInfo, id int) {
			Map(workers[id].IP, task.toRPC())
			wg.Done()
		}(mapTask, workerID)
		workerID = (workerID + 1) % nWorkers
	}
	wg.Wait()
	log.Trace("[Master] End Map task")
}

func (ms *Master) DistributeReduceTask() {
	// TODO: call workers's gRPC functions
	log.Trace("[Master] Start Reduce task")
	var wg sync.WaitGroup
	workerID := 0
	for _, reduceTask := range ms.ReduceTasks {
		wg.Add(1)
		go func(task ReduceTaskInfo, id int) {
			Reduce(ms.Workers[id].IP, task.toRPC())
			wg.Done()
		}(reduceTask, workerID)
		workerID = (workerID + 1) % ms.numReducer
	}
	wg.Wait()
	log.Trace("[Master] End Reduce task")
}

func (ms *Master) EndWorkers() {
	log.Trace("[Master] End Workers Start")
	for _, worker := range ms.Workers {
		End(worker.IP)
	}
	log.Trace("[Master] End Workers done")
}
