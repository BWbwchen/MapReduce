package master

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func (ms *Master) PeriodicHealthCheck() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		<-ticker.C
		ms.checkWorkersHealth()
	}

}

func (ms *Master) checkWorkersHealth() {
	var wg sync.WaitGroup
	for _, worker := range ms.Workers {
		wg.Add(1)
		go func(w WorkerInfo) {
			checkHealth(w)
			wg.Done()
		}(worker)
	}
	wg.Wait()
}

func checkHealth(worker WorkerInfo) {
	worker.WorkerState = Health(worker.IP)
	state := ""
	switch worker.WorkerState {
	case WORKER_IDLE:
		state = "Worker IDLE"
	case WORKER_BUSY:
		state = "Worker Busy"
	}
	log.Info("[Master] Worker state is :", state)
}
