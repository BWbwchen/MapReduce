package master

import (
	"net"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	IDLE int = iota
	INPROGRESS
	COMPLETED
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func StartMaster(files []string, nReduce int, addr string) {
	// start gRPC server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	ms := NewMaster(nReduce)
	baseServer := grpc.NewServer()
	rpc.RegisterMasterServer(baseServer, ms)
	go baseServer.Serve(listener)

	log.Info("[Master] Master gRPC server start")

	// Check the worker is enough
	ms.(*Master).WaitForEnoughWorker()

	// Split input file (100,000 lines per chunk)
	ms.(*Master).DistributeWork(files)

	ms.(*Master).DistributeMapTask()

	// TODO: uncomment
	ms.(*Master).DistributeReduceTask()

	ms.(*Master).EndWorkers()
}
