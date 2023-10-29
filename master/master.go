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
	crashChan    chan string
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
		crashChan:    make(chan string, 100),
	}
}

// gRPC functions

func (ms *Master) WorkerRegister(ctx context.Context, in *rpc.WorkerInfo) (*rpc.RegisterResult, error) {
	var num int
	ms.mux.Lock()
	ms.Workers = append(ms.Workers, newWorker(in.Uuid, in.Ip))
	ms.numWorkers++

	num = ms.numWorkers
	ms.mux.Unlock()
	log.Info("[Master] Worker register success")
	return &rpc.RegisterResult{Result: true, Id: int64(num - 1)}, nil
}

func (ms *Master) UpdateIMDInfo(ctx context.Context, in *rpc.IMDInfo) (*rpc.UpdateResult, error) {
	ms.mux.Lock()
	for i, f := range in.Filenames {
		ms.ReduceTasks[i].IMDs = append(ms.ReduceTasks[i].IMDs,
			IMDInfo{
				IP:       ms.serviceDiscovey(in.Uuid),
				FileName: f,
			})
	}
	ms.mux.Unlock()
	log.Info(fmt.Sprintf("[Master] %v update IMD info success", in.Uuid))
	return &rpc.UpdateResult{Result: true}, nil
}

func (ms *Master) serviceDiscovey(uuid string) string {
	var ip string

	for i := range ms.Workers {
		if ms.Workers[i].UUID == uuid {
			ip = ms.Workers[i].getIP()
		}
	}

	return ip
}

// Normal Functions

func (ms *Master) waitForEnoughWorker() {
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

func (ms *Master) distributeWork(files []string) {
	log.Trace("[Master] Start distribute workload")
	numWorkers := ms.totalWorkers
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

			ms.MapTasks[i].addFile(file, from, from+workLoad)
			from += workLoad
		}
	}
	log.Trace("[Master] End distribute workload")
}

func (ms *Master) availableWorkers(num int) ([]*WorkerInfo, int) {
	log.Info("[Master] Finding available workers to execute ", num, "/", len(ms.Workers))
	var retInfo []*WorkerInfo
	total := 0
	broken := 0

LOOP:
	for {
		// ms.mux.Lock()
		broken = 0
		for i := range ms.Workers {
			if total >= num {
				break LOOP
			}
			if ms.Workers[i].Health() {
				retInfo = append(retInfo, &ms.Workers[i])
				ms.Workers[i].SetState(WORKER_BUSY)
				total += 1
			} else if ms.Workers[i].Broken() {
				broken += 1
			}
		}
		if broken == len(ms.Workers) {
			log.Panic("No enough worker!", total, num)
		}
		// ms.mux.Unlock()
	}

	retWorkers := len(retInfo)
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

func (ms *Master) orDone(finish, crashChan <-chan string) <-chan string {
	count := 0
	valStream := make(chan string)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-finish:
				count += 1
				if count == len(ms.MapTasks) {
					return
				}
			case crashUUID, ok := <-crashChan:
				if !ok {
					if count != len(ms.MapTasks) {
						log.Panic("faild to distribute all MapTask")
						return
					}
				}

				valStream <- crashUUID

			}
		}
	}()
	return valStream
}

func (ms *Master) distributeMapTask() {
	log.Trace("[Master] Start Map task")

	taskStates := sync.Map{}

	finish := make(chan string, 100)
	defer close(finish)
	// crashChan := make(chan MapTaskInfo, 100)
	workerID := 0
	workers, nWorkers := ms.availableWorkers(ms.totalWorkers)
	for _, mapTask := range ms.MapTasks {
		mapTask.setState(TASK_INPROGRESS)
		go func(task MapTaskInfo, id int) {
			taskStates.Store(workers[id].UUID, task)
			done := Map(workers[id].IP, task.toRPC())
			if !done {
				// task.setState(TASK_IDLE)
				workers[id].SetState(WORKER_UNKNOWN)
				ms.crashChan <- workers[id].UUID
				// crashChan <- task
			} else {
				// task.setState(TASK_COMPLETED)
				workers[id].SetState(WORKER_IDLE)
				finish <- workers[id].UUID
			}
		}(mapTask, workerID)
		workerID = (workerID + 1) % nWorkers
	}

	for crashUUID := range ms.orDone(finish, ms.crashChan) {
		workers, _ := ms.availableWorkers(1)
		log.Info("[Master] Re-execute Map Task from ", crashUUID, " To ", workers[0].UUID)
		if crashUUID == workers[0].UUID {
			log.Panic("Self loop")
		}
		reExecuteTask, ok := taskStates.Load(crashUUID)
		if !ok {
			log.Panic("Load task states from crashUUID fail")
		}
		go func(task MapTaskInfo, id int) {
			// taskStates.Delete(crashUUID)
			taskStates.Store(workers[id].UUID, task)
			done := Map(workers[id].IP, task.toRPC())
			if !done {
				task.setState(TASK_IDLE)
				workers[id].SetState(WORKER_UNKNOWN)
				ms.crashChan <- workers[id].UUID
			} else {
				task.setState(TASK_COMPLETED)
				workers[id].SetState(WORKER_IDLE)
				finish <- workers[id].UUID
			}
		}(reExecuteTask.(MapTaskInfo), 0)
	}

	log.Trace("[Master] End Map task")
}

func (ms *Master) distributeReduceTask() {
	log.Trace("[Master] Start Reduce task")

	taskStates := sync.Map{}

	finish := make(chan string, 100)
	defer close(finish)
	workerID := 0
	workers, nWorkers := ms.availableWorkers(ms.numReducer)
	for _, reduceTask := range ms.ReduceTasks {
		reduceTask.SetState(TASK_INPROGRESS)
		go func(task ReduceTaskInfo, id int) {
			taskStates.Store(workers[id].UUID, task)
			workers[id].SetState(WORKER_BUSY)
			done := Reduce(workers[id].IP, task.toRPC())
			if !done {
				task.SetState(TASK_IDLE)
				workers[id].SetState(WORKER_UNKNOWN)
				ms.crashChan <- workers[id].UUID
			} else {
				task.SetState(TASK_COMPLETED)
				workers[id].SetState(WORKER_IDLE)
				finish <- workers[id].UUID
			}
		}(reduceTask, workerID)
		workerID = (workerID + 1) % nWorkers
	}

	for crashUUID := range ms.orDone(finish, ms.crashChan) {
		workers, _ = ms.availableWorkers(1)

		log.Info("[Master] Re-execute Map Task from ", crashUUID, "To ", workers[0].UUID)
		reExecuteTask, _ := taskStates.Load(crashUUID)
		go func(task ReduceTaskInfo, id int) {
			taskStates.Delete(crashUUID)
			taskStates.Store(workers[id].UUID, task)
			done := Reduce(workers[id].IP, task.toRPC())
			if !done {
				task.SetState(TASK_IDLE)
				workers[id].SetState(WORKER_UNKNOWN)
				ms.crashChan <- workers[id].UUID
			} else {
				task.SetState(TASK_COMPLETED)
				workers[id].SetState(WORKER_IDLE)
				finish <- workers[id].UUID
			}
		}(reExecuteTask.(ReduceTaskInfo), 0)
	}

	log.Trace("[Master] End Reduce task")
}

func (ms *Master) endWorkers() {
	log.Trace("[Master] End Workers Start")
	for i := range ms.Workers {
		End(ms.Workers[i].IP)
	}
	log.Trace("[Master] End Workers done")
}
