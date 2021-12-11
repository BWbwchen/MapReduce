package master

import (
	// log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// If there are 3 continuous unknow, we thought that that worker is dead.
// stop that work and make other deal with that.

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
		// worker.WorkerState = Health(worker.IP)
		wg.Add(1)
		go func(w WorkerInfo) {
			checkHealth(w)
			wg.Done()
		}(worker)
	}
	wg.Wait()
}

func checkHealth(worker WorkerInfo) {
	// worker.WorkerState = Health(worker.IP)
	// state := ""
	// switch worker.WorkerState {
	// case WORKER_IDLE:
	// 	state = "Worker IDLE"
	// case WORKER_BUSY:
	// 	state = "Worker Busy"
	// case WORKER_UNKNOWN:
	// 	state = "Worker Dead"
	// }

	// log.Info("[Master] Worker state is :", state)
}
