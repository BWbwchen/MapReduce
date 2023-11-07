package master

import (
	"context"
	"os"
	"testing"

	"github.com/BWbwchen/MapReduce/master/mocks"
	"github.com/BWbwchen/MapReduce/rpc"
)

func TestWorkerRegister(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}
	info := rpc.WorkerInfo{
		Uuid: "uuid",
		Ip:   "ip",
	}
	master.WorkerRegister(context.Background(), &info)

	if master.numWorkers != 1 {
		t.Error("number of workers is not correct!")
	}
	workerInfo := WorkerInfo{UUID: "uuid", IP: "ip", WorkerState: WORKER_IDLE}
	if master.Workers[0] != workerInfo {
		t.Error("worker info is not correct")
	}
}

func TestDistributeWork(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}

	fileNames, err := createTestFiles()
	if err != nil {
		t.Error(err)
	}

	master.distributeWork(fileNames)
	for counter, task := range master.MapTasks {
		if counter == 0 {
			for _, fileInfo := range task.Files {
				if fileInfo.From != 0 && fileInfo.To != 2 {
					t.Error("work didn`t distributed correctly")
				}
			}
		}
	}
}

func TestDistributeMapTask(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}

	workerInfo := WorkerInfo{UUID: "uuid", IP: "ip", WorkerState: WORKER_IDLE}
	workerInfo1 := WorkerInfo{UUID: "uuid1", IP: "ip1", WorkerState: WORKER_IDLE}

	mapTaskFile := "test"

	master.MapTasks = []MapTaskInfo{{}}
	master.MapTasks[0].addFile(mapTaskFile, 0, 1)
	master.Workers = append(master.Workers, workerInfo, workerInfo1)

	mocks.Result = true
	master.distributeMapTask()
	req, _ := mocks.Request.(*rpc.MapInfo)
	if req.Files[0].FileName != mapTaskFile || req.Files[0].From != 0 || req.Files[0].To != 1 {
		t.Error("request to worker is not correct")
	}
}

func TestDistributeMapTaskFailsTolerant(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}

	workerInfo := WorkerInfo{UUID: "uuid", IP: "ip", WorkerState: WORKER_IDLE}
	workerInfo1 := WorkerInfo{UUID: "uuid1", IP: "ip1", WorkerState: WORKER_IDLE}

	mapTaskFile := "test"

	master.MapTasks = []MapTaskInfo{{}}
	master.MapTasks[0].addFile(mapTaskFile, 0, 1)
	master.Workers = append(master.Workers, workerInfo, workerInfo1)

	mocks.Result = false
	master.distributeMapTask()
	req, _ := mocks.Request.(*rpc.MapInfo)
	if req.Files[0].FileName != mapTaskFile || req.Files[0].From != 0 || req.Files[0].To != 1 {
		t.Error("request to worker is not correct")
	}
}

func TestDistributeReduceTask(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}

	workerInfo := WorkerInfo{UUID: "uuid", IP: "ip", WorkerState: WORKER_IDLE}
	workerInfo1 := WorkerInfo{UUID: "uuid1", IP: "ip1", WorkerState: WORKER_IDLE}

	IMDFileName := "testReduceTasks"

	master.ReduceTasks = []ReduceTaskInfo{{IMDs: []IMDInfo{
		{FileName: IMDFileName},
	}}}

	master.Workers = append(master.Workers, workerInfo, workerInfo1)

	mocks.Result = true
	master.distributeReduceTask()
	req, _ := mocks.Request.(*rpc.ReduceInfo)
	if req.Files[0].Filename != IMDFileName {
		t.Error("request to worker is not correct")
	}
}

func TestDistributeReduceTaskFailsTolerant(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}

	workerInfo := WorkerInfo{UUID: "uuid", IP: "ip", WorkerState: WORKER_IDLE}
	workerInfo1 := WorkerInfo{UUID: "uuid1", IP: "ip1", WorkerState: WORKER_IDLE}

	IMDFileName := "testReduceTasks"

	master.ReduceTasks = []ReduceTaskInfo{{IMDs: []IMDInfo{
		{FileName: IMDFileName},
	}}}

	master.Workers = append(master.Workers, workerInfo, workerInfo1)

	mocks.Result = false
	master.distributeReduceTask()
	req, _ := mocks.Request.(*rpc.ReduceInfo)
	if req.Files[0].Filename != IMDFileName {
		t.Error("request to worker is not correct")
	}
}

func TestEndWorkers(t *testing.T) {
	master := NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}
	master.endWorkers()
}

func createTestFiles() ([]string, error) {
	fileNames := []string{}
	file1, err := os.CreateTemp("", "tmp-file")
	if err != nil {
		return nil, err
	}
	defer file1.Close()
	fileNames = append(fileNames, file1.Name())

	_, err = file1.WriteString("this is some text for testing \n next line \n nextline")
	if err != nil {
		return nil, err
	}

	file2, err := os.CreateTemp("", "tmp-file")
	if err != nil {
		return nil, err
	}
	defer file2.Close()
	fileNames = append(fileNames, file2.Name())

	_, err = file2.WriteString("this is some text for testing\n next line \n nextline")
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}
