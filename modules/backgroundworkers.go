// Package modules contains the implementation of the background workers service.
//
// The background workers service is responsible for running background workers
// as goroutines, and for providing a way to register and start workers.
//
// The service is useful for running background tasks, such as sending emails,
// in a separate goroutine from the main application goroutine.
package modules

import (
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
)

// background_worker_svc is a service that runs background workers as goroutines.
type background_worker_svc struct {
	Logger  logger.ILogger
	Config  config.Config
	Workers []Worker
}

// Start starts the background workers.
func (b *background_worker_svc) Start() {
	for _, worker := range b.Workers {
		go func(worker Worker) {
			// Create a ticker that will fire at the specified interval.
			ticker := time.NewTicker(worker.Interval)
			// Run the task for each tick.
			for range ticker.C {
				worker.Task()
			}
		}(worker)
	}
}
