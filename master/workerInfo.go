package master

import "sync"

type WorkerInfo struct {
	UUID        string
	IP          string
	WorkerState int
	mux         sync.Mutex
}

func newWorker(uuid string, ip string) WorkerInfo {
	return WorkerInfo{
		UUID:        uuid,
		IP:          ip,
		WorkerState: WORKER_IDLE,
	}
}

func (w *WorkerInfo) getIP() string {
	w.mux.Lock()
	ret := w.IP
	w.mux.Unlock()
	return ret
}

func (w *WorkerInfo) Health() bool {
	w.mux.Lock()
	ret := w.WorkerState == WORKER_IDLE
	w.mux.Unlock()
	return ret
}

func (w *WorkerInfo) SetState(state int) {
	w.mux.Lock()
	w.WorkerState = state
	w.mux.Unlock()
}

func (w *WorkerInfo) Broken() bool {
	return w.WorkerState == WORKER_UNKNOWN
}
