package master

import (
	"context"
	"time"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Map(workerIP string, m *rpc.MapInfo) bool {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return false
	}
	defer conn.Close()
	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Map(ctx, m)
	if err != nil {
		log.Warn("[Master]: " + err.Error())
		return false
	}
	return r.Result
}

func Reduce(workerIP string, m *rpc.ReduceInfo) bool {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return false
	}
	defer conn.Close()
	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.Reduce(ctx, m)
	if err != nil {
		log.Warn("[Master]: " + err.Error())
		return false
	}
	return r.Result
}

func End(workerIP string) bool {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return false
	}
	defer conn.Close()
	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c.End(ctx, &rpc.Empty{})

	return true
}

func Health(workerIP string) int {
	conn, err := grpc.Dial(workerIP, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
		return int(WORKER_UNKNOWN)
	}
	defer conn.Close()
	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	r, err := c.Health(ctx, &rpc.Empty{})
	if err != nil {
		return int(WORKER_UNKNOWN)
	}

	log.Info("[Master] Worker State ", int(r.State))
	return int(r.State)
}
