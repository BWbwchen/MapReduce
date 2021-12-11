package master

import (
	"github.com/BWbwchen/MapReduce/rpc"
	"github.com/google/uuid"
)

type ReduceTaskInfo struct {
	IMDs      []IMDInfo
	TaskState int
	UUID      string
}

type IMDInfo struct {
	IP       string
	FileName string
}

func newReduceTask() ReduceTaskInfo {
	return ReduceTaskInfo{
		UUID:      uuid.New().String(),
		TaskState: TASK_IDLE,
	}
}

func newReduceTasks(size int) []ReduceTaskInfo {
	var ret []ReduceTaskInfo
	for i := 0; i < size; i++ {
		ret = append(ret, newReduceTask())
	}
	return ret
}

func (mt *ReduceTaskInfo) toRPC() *rpc.ReduceInfo {
	ret := &rpc.ReduceInfo{}

	for _, finfo := range mt.IMDs {
		ret.Files = append(ret.Files, &rpc.ReduceFileInfo{
			Ip:       finfo.IP,
			Filename: finfo.FileName,
		})
	}

	return ret
}

func (mt *ReduceTaskInfo) SetState(state int) {
	mt.TaskState = state
}
