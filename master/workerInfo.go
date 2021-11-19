package master

type WorkerInfo struct {
	UUID        string
	IP          string
	WorkerState int
}

func NewWorker(uuid string, ip string) WorkerInfo {
	return WorkerInfo{
		UUID: uuid,
		IP:   ip,
	}
}

func (w *WorkerInfo) GetIP() string {
	return w.IP
}
