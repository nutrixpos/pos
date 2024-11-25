package modules

import (
	"time"

	"github.com/elmawardy/nutrix/backend/common/config"
	"github.com/elmawardy/nutrix/backend/common/logger"
)

type background_worker_svc struct {
	Logger  logger.ILogger
	Config  config.Config
	Workers []Worker
}

func (b *background_worker_svc) Start() {
	for _, worker := range b.Workers {
		go func(worker Worker) {
			ticker := time.NewTicker(worker.Interval)
			for range ticker.C {
				worker.Task()
			}
		}(worker)
	}
}
