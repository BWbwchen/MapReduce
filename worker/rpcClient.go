package worker

import (
	"context"
	"time"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func WorkerRegister(w *rpc.WorkerInfo) int {
	conn, err := grpc.Dial(MasterIP, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := rpc.NewMasterClient(conn)
	log.Trace("New Master Client")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	log.Trace("With Time out")
	defer cancel()

	log.Trace("Start RPC call")
	r, err := c.WorkerRegister(ctx, w)
	log.Trace("End RPC call")
	if err != nil {
		panic(err)
	}

	if r.Result == false {
		panic("Register Error")
	}

	return int(r.Id)
}

func UpdateIMDInfo(u *rpc.IMDInfo) bool {
	conn, err := grpc.Dial(MasterIP, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	c := rpc.NewMasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UpdateIMDInfo(ctx, u)
	if err != nil {
		log.Panic(err)
	}

	if r.Result == false {
		log.Panic("Update IMD Info Error")
	}

	return r.Result
}

func GetIMDData(ip string, filename string) []KV {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetIMDData(ctx, &rpc.IMDLoc{
		Filename: filename,
	})

	var kvs []KV
	for _, kv := range r.Kvs {
		kvs = append(kvs, newKV(kv.Key, kv.Value))
	}

	return kvs
}
