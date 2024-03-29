package master

import (
	"context"
	"time"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type RpcClient interface {
	Connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient)
	Map(workerIP string, m *rpc.MapInfo) bool
	Reduce(workerIP string, m *rpc.ReduceInfo) bool
	End(workerIP string) bool
	Health(workerIP string) int
}

type workerClient struct{}

func (client *workerClient) Connect(workerIP string) (*grpc.ClientConn, rpc.WorkerClient) {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return nil, nil
	}

	return conn, rpc.NewWorkerClient(conn)
}

func (client *workerClient) Map(workerIP string, m *rpc.MapInfo) bool {
	conn, c := client.Connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Map(ctx, m)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes()
			log.Warn("[Master]: " + respErr.Message())
		} else {
			log.Warn("[Master]: " + err.Error())
		}
		return false
	}
	return r.Result
}

func (client *workerClient) Reduce(workerIP string, m *rpc.ReduceInfo) bool {
	conn, c := client.Connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Reduce(ctx, m)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes(respErr.Code())
			log.Warn("[Master]: " + respErr.Message())
		} else {
			log.Warn("[Master]: " + err.Error())
		}
		return false
	}
	return r.Result
}

func (client *workerClient) End(workerIP string) bool {
	conn, c := client.Connect(workerIP)
	if conn == nil {
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.End(ctx, &rpc.Empty{})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes(respErr.Code())
			log.Warn("[Master]: " + respErr.Message())
		} else {
			log.Warn("[Master]: " + err.Error())
		}
		return false
	}

	return true
}

func (client *workerClient) Health(workerIP string) int {
	conn, c := client.Connect(workerIP)
	if conn == nil {
		return int(WORKER_UNKNOWN)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	r, err := c.Health(ctx, &rpc.Empty{})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes(respErr.Code())
			log.Warn("[Master]: " + respErr.Message())
		} else {
			log.Warn("[Master]: " + err.Error())
		}
		return int(WORKER_UNKNOWN)
	}

	log.Info("[Master] Worker State ", int(r.State))
	return int(r.State)
}
