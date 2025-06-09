package hubsync

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/modules"
	auth_mw "github.com/nutrixpos/pos/modules/auth/middlewares"
	"github.com/nutrixpos/pos/modules/core/middlewares"
	"github.com/nutrixpos/pos/modules/hubsync/handlers"
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

func (c *HubSyncModule) RegisterHttpHandlers(router *mux.Router, prefix string) {

	auth_svc := auth_mw.NewZitadelAuth(c.Config)

	c.Logger.Info("Successfully conntected to Zitadel")

	router.Handle(prefix+"/api/settings", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.PatchSettings(c.Config, c.Logger), "admin"))).Methods("PATCH", "OPTIONS")
	router.Handle(prefix+"/api/settings", middlewares.AllowCors(auth_svc.AllowAnyOfRoles(handlers.GetSettings(c.Config, c.Logger), "admin"))).Methods("GET", "OPTIONS")
}

// RegisterBackgroundWorkers registers background workers.
func (c *HubSyncModule) RegisterBackgroundWorkers() []modules.Worker {

	workers := make([]modules.Worker, 0)

	settingsSvc := services.SettingsSvc{
		Config: c.Config,
		Logger: c.Logger,
	}

	seconds := 60

	settings, err := settingsSvc.Get()
	if err != nil {
		c.Logger.Error(fmt.Sprintf("Failed to get settings at hubsync.RegisterBackgroundWorkers: %s", err.Error()))
	}

	if err == nil {

		if settings.Settings.SyncInterval > 0 {
			seconds = int(settings.Settings.SyncInterval)
		} else {
			c.Logger.Info("Sync interval is less than or equal 0, defaulting to 60 seconds")
		}

		c.Settings = settings
	}

	workers = append(workers, modules.Worker{
		Interval: time.Duration(seconds) * time.Second,
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
