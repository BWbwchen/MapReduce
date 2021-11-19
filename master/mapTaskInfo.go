package master

import (
	"github.com/BWbwchen/MapReduce/rpc"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MapTaskInfo struct {
	Files     []FileInfo
	TaskState int
	UUID      string
}

type FileInfo struct {
	FileName string
	From     int
	To       int
}

func (mt *MapTaskInfo) AddFile(file string, from int, to int) {
	mt.Files = append(mt.Files, FileInfo{
		FileName: file,
		From:     from,
		To:       to,
	})
	log.Trace("Add map input file")
}

func newMapTask() MapTaskInfo {
	return MapTaskInfo{
		UUID:      uuid.New().String(),
		TaskState: IDLE,
	}
}

func newMapTasks(size int) []MapTaskInfo {
	var ret []MapTaskInfo
	for i := 0; i < size; i++ {
		ret = append(ret, newMapTask())
	}
	return ret
}

func (mt *MapTaskInfo) toRPC() *rpc.MapInfo {
	ret := &rpc.MapInfo{}

	for _, finfo := range mt.Files {
		ret.Files = append(ret.Files, &rpc.MapFileInfo{
			FileName: finfo.FileName,
			From:     int64(finfo.From),
			To:       int64(finfo.To),
		})
	}

	return ret
}
