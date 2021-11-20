package worker

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

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
	rpc.UnimplementedWorkerServer
}

func newWorker(nReduce int, inRAM bool) rpc.WorkerServer {
	return &Worker{
		UUID:       uuid.New().String(),
		nReduce:    nReduce,
		Chan:       newMrContext(),
		EndChan:    make(chan bool),
		storeInRAM: inRAM,
	}
}

// gRPC functions

func (wr *Worker) Map(ctx context.Context, in *rpc.MapInfo) (*rpc.Result, error) {
	log.Info("[Worker] Start Map")
	log.Trace("[Worker] Start Mapping")
	done := make(chan int, 100)
	for _, fInfo := range in.Files {
		content := partialContent(fInfo)
		// Call wr.Mapf
		go func(f0 *rpc.MapFileInfo, c0 string) {
			wr.Mapf(f0.FileName, c0, wr.Chan)
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
		case mapKV, haveKV := <-wr.Chan.MapChan:
			if haveKV {
				imdKV[i] = append(imdKV[i], mapKV)
				i = (i + 1) % wr.nReduce
			} else {
				break LOOP
			}

		case <-done:
			count++
			if count == len(in.Files) {
				close(wr.Chan.MapChan)
			}
		}

	}
	log.Trace("[Worker] End partition intermediate kv")

	log.Trace("[Worker] Write intermediate kv to file")
	filenames := writeIMDToLocalFile(imdKV, wr.UUID, wr.storeInRAM)
	log.Trace("[Worker] End Write intermediate kv to file")

	log.Trace("[Worker] Tell Master the intermediate info")
	// Return to the Master
	UpdateIMDInfo(&rpc.IMDInfo{
		Uuid:      wr.UUID,
		Filenames: filenames,
	})
	log.Trace("[Worker] Finish Tell Master the intermediate info")
	log.Info("[Worker] Finish Map Task")

	return &rpc.Result{Result: true}, nil
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
	content := ""
	for _, kv := range kvs {
		content += fmt.Sprintf("%v %v\n", kv.Key, kv.Value)
	}

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
	log.Trace("[Worker] Get intermediate file")
	var imdKVs []KV
	for _, fInfo := range in.Files {
		imdKVs = append(imdKVs, GetIMDData(fInfo.Ip, fInfo.Filename)...)
	}

	log.Trace("[Worker] Sort intermediate KV")
	// Sort
	sort.Sort(byKey(imdKVs))

	// outputFile := fmt.Sprintf("output/mr-out-%v.txt", wr.UUID)
	outputFile := fmt.Sprintf("output/mr-out-%v.txt", wr.ID)
	ofile, _ := os.Create(outputFile)

	log.Trace("[Worker] Start Reducing")
	// Reduce all the intermediate KV
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
		wr.Reducef(imdKVs[i].Key, values, wr.Chan)

		output := <-wr.Chan.ReduceChan
		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", output.Key, output.Value)

		i = j
	}
	log.Trace("[Worker] End Reducing")
	log.Info("[Worker] End Reduce")

	return &rpc.Result{Result: true}, nil
}

func (wr *Worker) GetIMDData(ctx context.Context, in *rpc.IMDLoc) (*rpc.KVs, error) {
	log.Info("[Worker] RPC Get intermediate file")
	return &rpc.KVs{
		Kvs: generateIMDKV(in.Filename),
	}, nil
}

func (wr *Worker) End(ctx context.Context, in *rpc.Empty) (*rpc.Empty, error) {
	log.Info("[Worker] End worker")
	wr.EndChan <- true
	return &rpc.Empty{}, nil
}

// Normal Functions

func (wr *Worker) setID(id int) {
	wr.ID = id
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

func generateIMDKV(file string) []*rpc.KV {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var Kvs []*rpc.KV

	for scanner.Scan() {
		var kv rpc.KV
		line := strings.Split(scanner.Text(), " ")
		kv.Key = line[0]
		kv.Value = line[1]
		Kvs = append(Kvs, &kv)
	}

	return Kvs
}
