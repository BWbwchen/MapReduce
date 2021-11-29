package master

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (ms *Master) periodicHealthCheck() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		<-ticker.C
		ms.checkWorkersHealth()
	}

}

func (ms *Master) checkWorkersHealth() {
	for _, worker := range ms.Workers {
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
}
