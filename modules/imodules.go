// Package modules provides the main logic for the module system in the application.
// It allows registering modules, starting them, and stopping them.
// It also provides a way to register HTTP handlers for the modules.
package modules

import (
	"time"

	"github.com/gorilla/mux"
)

// saved_module_builders is a map of module builders that have been saved.
// These are loaded by the app manager when the app starts.
var saved_module_builders = make(map[string]*ModuleBuilder)

// IBaseModule is the interface that all modules must implement.
// It provides methods for starting and stopping the module.
type IBaseModule interface {
	// OnStart is called when the module is started.
	OnStart() func() error
	// OnEnd is called when the module is stopped.
	OnEnd() func()
}

// Worker represents a background worker.
// Interval is the duration between each run of the task.
// Task is the function that is called every Interval.
type Worker struct {
	Interval time.Duration
	Task     func()
}

// IHttpModule is an interface that modules can implement to add HTTP handlers.
// RegisterHttpHandlers is called by the app manager to register the HTTP handlers.
type IHttpModule interface {
	RegisterHttpHandlers(r *mux.Router, prefix string)
}

// IBackgroundWorkerModule is an interface that modules can implement to add background workers.
// RegisterBackgroundWorkers is called by the app manager to register the background workers.
type IBackgroundWorkerModule interface {
	RegisterBackgroundWorkers() []Worker
}

// ISeederModule is an interface that modules can implement to add seeders.
// Seed is called by the app manager to seed the database.
// GetSeedables is called by the app manager to get the list of seedables.
type ISeederModule interface {
	Seed(entities []string, is_new_only bool) error
	GetSeedables() (entities []string, err error)
}
