package worker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/BWbwchen/MapReduce/rpc"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MapFormat (func(string, string, MrContext))
type ReduceFormat (func(string, []string, MrContext))

type Worker struct {
	UUID       string
	ID         int
	nReduce    int
	Mapf       MapFormat
	Reducef    ReduceFormat
	Chan       MrContext
	EndChan    chan bool
	storeInRAM bool
	State      rpc.WorkerState_State
	Client     RpcClient
	mux        sync.Mutex
	rpc.UnimplementedWorkerServer
}

func newWorker(nReduce int, inRAM bool) rpc.WorkerServer {
	conn, master := Connect()
	return &Worker{
		UUID:       uuid.New().String(),
		nReduce:    nReduce,
		Chan:       newMrContext(),
		EndChan:    make(chan bool),
		Client:     &masterClient{master: master, conn: conn},
		storeInRAM: inRAM,
		State:      rpc.WorkerState_IDLE,
	}
}

// gRPC functions

func (wr *Worker) Map(ctx context.Context, in *rpc.MapInfo) (*rpc.Result, error) {
	log.Info("[Worker] Start Map")

	wr.setWorkerState(rpc.WorkerState_BUSY)

	log.Trace("[Worker] Start Mapping")
	done := make(chan int, 100)
	mapChan := newMrContext()
	for _, fInfo := range in.Files {
		content := partialContent(fInfo)
		go func(f0 *rpc.MapFileInfo, c0 string) {
			wr.Mapf(f0.FileName, c0, mapChan)
			done <- 1
		}(fInfo, content)
	}
	log.Trace("[Worker] Finish Mapping")

	imdKV := make([][]KV, wr.nReduce)

	// Get intermediate KV
	// Partition result into R piece
	log.Trace("[Worker] Start partition intermediate kv")
	i := 0
	count := 0

LOOP:
	for {
		select {
		case mapKV, haveKV := <-mapChan.Chan:
			if haveKV {
				imdKV[i] = append(imdKV[i], mapKV)
				i = (i + 1) % wr.nReduce
			} else {
				break LOOP
			}

		case <-done:
			count++
			if count == len(in.Files) {
				close(mapChan.Chan)
			}
		}

	}
	log.Trace("[Worker] End partition intermediate kv")

	log.Trace("[Worker] Write intermediate kv to file")
	filenames := writeIMDToLocalFile(imdKV, wr.UUID, wr.storeInRAM)
	log.Trace("[Worker] End Write intermediate kv to file")

	log.Trace("[Worker] Tell Master the intermediate info")
	// Return to the Master
	wr.Client.UpdateIMDInfo(&rpc.IMDInfo{
		Uuid:      wr.UUID,
		Filenames: filenames,
	})
	log.Trace("[Worker] Finish Tell Master the intermediate info")
	log.Info("[Worker] Finish Map Task")
	wr.setWorkerState(rpc.WorkerState_IDLE)

	return &rpc.Result{Result: true}, nil
}

func partialContent(fInfo *rpc.MapFileInfo) string {
	f, err := os.Open(fInfo.FileName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	line := 0
	content := ""

	for scanner.Scan() {
		if int(fInfo.From) <= line && line < int(fInfo.To) {
			content += scanner.Text() + "\n"
		} else if line >= int(fInfo.To) {
			break
		}
		line++
	}

	return content
}

func writeIMDToLocalFile(imdKV [][]KV, uuid string, inRAM bool) []string {
	var filenames []string
	// Write to local file
	filenameChan := make(chan string, 20)
	finishChan := make(chan bool, 2)
	for taskId, kvs := range imdKV {
		go func(t int, s []KV) {
			writeIMDToLocalFileParallel(t, s, uuid, inRAM, filenameChan, finishChan)
		}(taskId, kvs)
	}

	count := 0
LOOP:
	for {
		select {
		case f, more := <-filenameChan:
			if more {
				filenames = append(filenames, f)
			} else {
				break LOOP
			}
		case <-finishChan:
			count++
			if count == len(imdKV) {
				close(filenameChan)
			}
		}
	}
	return filenames
}

func writeIMDToLocalFileParallel(taskId int, kvs []KV, uuid string, inRAM bool, output chan string, finish chan bool) {
	content_byte, _ := json.Marshal(kvs)
	content := string(content_byte)

	var fname string
	if inRAM {
		fname = fmt.Sprintf("/dev/shm/imd-%v-%v.txt", uuid, taskId)
	} else {
		fname = fmt.Sprintf("output/imd-%v-%v.txt", uuid, taskId)
	}
	output <- fname
	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(content)
	finish <- true
}

func (wr *Worker) Reduce(ctx context.Context, in *rpc.ReduceInfo) (*rpc.Result, error) {
	log.Info("[Worker] Start Reduce")

	wr.setWorkerState(rpc.WorkerState_BUSY)

	log.Trace("[Worker] Get intermediate file")
	var imdKVs []KV
	for _, fInfo := range in.Files {
		imdKVs = append(imdKVs, wr.Client.GetIMDData(fInfo.Ip, fInfo.Filename)...)
	}

	log.Trace("[Worker] Sort intermediate KV")
	// Sort
	sort.Sort(byKey(imdKVs))

	outputFile := fmt.Sprintf("mr-out-%v.txt", wr.ID)
	ofile, _ := os.Create(outputFile)

	log.Trace("[Worker] Start Reducing")
	// Reduce all the intermediate KV
	reduceChan := newMrContext()
	i := 0
	for i < len(imdKVs) {
		j := i + 1
		for j < len(imdKVs) && imdKVs[j].Key == imdKVs[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, imdKVs[k].Value)
		}
		wr.Reducef(imdKVs[i].Key, values, reduceChan)

		output := <-reduceChan.Chan

		fmt.Fprintf(ofile, "%v %v\n", output.Key, output.Value)

		i = j
	}
	log.Trace("[Worker] End Reducing")
	log.Info("[Worker] End Reduce")
	wr.setWorkerState(rpc.WorkerState_IDLE)

	return &rpc.Result{Result: true}, nil
}

func (wr *Worker) GetIMDData(ctx context.Context, in *rpc.IMDLoc) (*rpc.JSONKVs, error) {
	log.Info("[Worker] RPC Get intermediate file")
	return &rpc.JSONKVs{
		Kvs: generateIMDKV(in.Filename),
	}, nil
}

func generateIMDKV(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	content := string(b)

	return content
}

func (wr *Worker) End(ctx context.Context, in *rpc.Empty) (*rpc.Empty, error) {
	log.Info("[Worker] End worker")
	wr.EndChan <- true
	return &rpc.Empty{}, nil
}

func (wr *Worker) setID(id int) {
	wr.ID = id
}

func (wr *Worker) Health(ctx context.Context, in *rpc.Empty) (*rpc.WorkerState, error) {
	log.Trace("[Worker] Health Check")

	wr.mux.Lock()
	state := wr.State
	wr.mux.Unlock()

	return &rpc.WorkerState{State: state}, nil
}

func (wr *Worker) setWorkerState(state rpc.WorkerState_State) {
	wr.mux.Lock()
	wr.State = state
	wr.mux.Unlock()
}
