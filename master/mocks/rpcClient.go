package mocks

import (
	"github.com/BWbwchen/MapReduce/rpc"
	"google.golang.org/grpc"
)

type WorkerClient struct{}

func (WorkerClient) Connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient) {
	conn := grpc.ClientConn{}

	return &conn, rpc.NewWorkerClient(&conn)
}

func (WorkerClient) Map(workerIP string, m *rpc.MapInfo) bool {
	return true
}

func (WorkerClient) Reduce(workerIP string, m *rpc.ReduceInfo) bool {
	return true
}

func (WorkerClient) End(workerIP string) bool {
	return true
}

func (WorkerClient) Health(workerIP string) int {
	return 0
}
