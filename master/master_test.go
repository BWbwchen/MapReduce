package master

import (
	"context"
	"os"
	"testing"

	"github.com/BWbwchen/MapReduce/master/mocks"
	"github.com/BWbwchen/MapReduce/rpc"
)

var master *Master

func init() {
	master = NewMaster(2, 1).(*Master)
	master.client = mocks.WorkerClient{}
}

func TestDistributeWork(t *testing.T) {

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
	info := rpc.WorkerInfo{
		Uuid: "uuid",
		Ip:   "ip",
	}
	info1 := rpc.WorkerInfo{
		Uuid: "uuid1",
		Ip:   "ip1",
	}

	master.WorkerRegister(context.Background(), &info)
	master.WorkerRegister(context.Background(), &info1)

	master.distributeMapTask()

}

func TestDistributeReduceTask(t *testing.T) {
	info := rpc.WorkerInfo{
		Uuid: "uuid",
		Ip:   "ip",
	}
	info1 := rpc.WorkerInfo{
		Uuid: "uuid1",
		Ip:   "ip1",
	}

	master.WorkerRegister(context.Background(), &info)
	master.WorkerRegister(context.Background(), &info1)

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
