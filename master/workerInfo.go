package master

type WorkerInfo struct {
	UUID        string
	IP          string
	WorkerState int
}

func newWorker(uuid string, ip string) WorkerInfo {
	return WorkerInfo{
		UUID:        uuid,
		IP:          ip,
		WorkerState: WORKER_IDLE,
	}
}

func (w *WorkerInfo) getIP() string {
	ret := w.IP
	return ret
}

func (w *WorkerInfo) Health() bool {
	ret := w.WorkerState == WORKER_IDLE
	return ret
}

func (w *WorkerInfo) SetState(state int) {
	w.WorkerState = state
}

func (w *WorkerInfo) Broken() bool {
	return w.WorkerState == WORKER_UNKNOWN
}
