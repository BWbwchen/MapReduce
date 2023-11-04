package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/BWbwchen/MapReduce/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type RpcClient interface {
	//Connect()
	WorkerRegister(w *rpc.WorkerInfo) int
	UpdateIMDInfo(u *rpc.IMDInfo) bool
	GetIMDData(ip string, filename string) []KV
}

type masterClient struct {
	master rpc.MasterClient
	conn   *grpc.ClientConn
}

func Connect(ip string) (*grpc.ClientConn, rpc.MasterClient) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.Warn(err)
	}
	return conn, rpc.NewMasterClient(conn)
}

func (client *masterClient) WorkerRegister(w *rpc.WorkerInfo) int {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	log.Trace("With Time out")
	defer cancel()

	log.Trace("Start RPC call")
	r, err := client.master.WorkerRegister(ctx, w)
	log.Trace("End RPC call")
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes()
			log.Panic(respErr.Message())
		} else {
			log.Panic(err)
		}
	}

	if !r.Result {
		panic("Register Error")
	}

	return int(r.Id)
}

func (client *masterClient) UpdateIMDInfo(u *rpc.IMDInfo) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.master.UpdateIMDInfo(ctx, u)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes()
			log.Panic(respErr.Message())
		} else {
			log.Panic(err)
		}
		return false
	}

	if !r.Result {
		log.Panic("Update IMD Info Error")
	}

	return r.Result
}

func (client *masterClient) GetIMDData(ip string, filename string) []KV {
	conn, _ := Connect(ip)

	c := rpc.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetIMDData(ctx, &rpc.IMDLoc{
		Filename: filename,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC
			//todo: we can improve error handling by using grpc statusCodes()
			log.Panic(respErr.Message())
		} else {
			log.Panic(err)
		}
	}

	var ret []KV

	err = json.Unmarshal([]byte(r.Kvs), &ret)
	if err != nil {
		panic(err)
	}

	return ret
}
