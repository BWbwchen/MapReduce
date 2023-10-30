package master

import (
	"context"
	"time"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type rpcClient interface {
	connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient)
	Map(workerIP string, m *rpc.MapInfo) bool
	Reduce(workerIP string, m *rpc.ReduceInfo) bool
	End(workerIP string) bool
	Health(workerIP string) int
}

type workerClient struct{}

func (client *workerClient) connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient) {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return nil, nil
	}

	c := rpc.NewWorkerClient(conn)

	return conn, c
}

func (client *workerClient) Map(workerIP string, m *rpc.MapInfo) bool {
	conn, c := client.connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Map(ctx, m)
	if err != nil {
		log.Warn("[Master]: " + err.Error())
		return false
	}
	return r.Result
}

func (client *workerClient) Reduce(workerIP string, m *rpc.ReduceInfo) bool {
	conn, c := client.connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Reduce(ctx, m)
	if err != nil {
		log.Warn("[Master]: " + err.Error())
		return false
	}
	return r.Result
}

func (client *workerClient) End(workerIP string) bool {
	conn, c := client.connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.End(ctx, &rpc.Empty{})

	return true
}

func (client *workerClient) Health(workerIP string) int {
	conn, c := client.connect(workerIP)
	if conn == nil {
		return int(WORKER_UNKNOWN)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	r, err := c.Health(ctx, &rpc.Empty{})
	if err != nil {
		return int(WORKER_UNKNOWN)
	}

	log.Info("[Master] Worker State ", int(r.State))
	return int(r.State)
}
