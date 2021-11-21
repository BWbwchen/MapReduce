package master

type WorkerInfo struct {
	UUID        string
	IP          string
	WorkerState int
}

func newWorker(uuid string, ip string) WorkerInfo {
	return WorkerInfo{
		UUID: uuid,
		IP:   ip,
	}
}

func (w *WorkerInfo) getIP() string {
	return w.IP
}
