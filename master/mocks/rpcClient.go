package mocks

import (
	"github.com/BWbwchen/MapReduce/rpc"
	"google.golang.org/grpc"
)

type WorkerClient struct{}

var Request interface{}
var Counter int
var Result bool

func (WorkerClient) Connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient) {
	conn := grpc.ClientConn{}

	return &conn, rpc.NewWorkerClient(&conn)
}

func (WorkerClient) Map(workerIP string, m *rpc.MapInfo) bool {
	Counter++
	Request = m
	// return true in second try for FailsTolerant
	if Counter > 1 {
		Result = true
	}
	return Result
}

func (WorkerClient) Reduce(workerIP string, m *rpc.ReduceInfo) bool {
	Counter++
	Request = m
	// return true in second try for FailsTolerant
	if Counter > 1 {
		Result = true
	}
	return Result
}

func (WorkerClient) End(workerIP string) bool {
	return Result
}

func (WorkerClient) Health(workerIP string) int {
	return 0
}
