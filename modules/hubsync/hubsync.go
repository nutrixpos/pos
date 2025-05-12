package hubsync

import (
	"fmt"
	"time"

	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules"
	"github.com/nutrixpos/pos/modules/hubsync/models"
	"github.com/nutrixpos/pos/modules/hubsync/services"
)

type HubSyncModule struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings models.Hubsync
}

// OnStart is called when the core module is started.
func (hs *HubSyncModule) OnStart() func() error {
	return func() error {

		err := hs.EnsureSeeded()
		if err != nil {
			return err
		}

		return nil
	}
}

// OnEnd is called when the core module is ended.
func (hs *HubSyncModule) OnEnd() func() {
	return func() {

	}
}

// RegisterBackgroundWorkers registers background workers.
func (c *HubSyncModule) RegisterBackgroundWorkers() []modules.Worker {

	workers := make([]modules.Worker, 0)

	workers = append(workers, modules.Worker{
		Interval: 1 * time.Minute,
		Task: func() {
			syncerSvc := services.SyncerService{
				Config: c.Config,
				Logger: c.Logger,
			}

			err := syncerSvc.Sync()
			if err != nil {
				c.Logger.Error(fmt.Sprintf("Error during sync: %s", err.Error()))
			}
		},
	})

	return workers
}

func (c *HubSyncModule) EnsureSeeded() error {

	seeder_svc := services.SeederService{
		Config: c.Config,
		Logger: c.Logger,
	}

	err := seeder_svc.Seed()
	if err != nil {
		return err
	}

	return nil
}
