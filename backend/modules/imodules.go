package modules

import (
	"time"

	"github.com/gorilla/mux"
)

var saved_module_builders = make(map[string]*ModuleBuilder)

type IBaseModule interface {
	OnStart() func()
	OnEnd() func()
}

type Worker struct {
	Interval time.Duration
	Task     func()
}

type IHttpModule interface {
	RegisterHttpHandlers(r *mux.Router, prefix string)
}

type IBackgroundWorkerModule interface {
	RegisterBackgroundWorkers() []Worker
}

type ISeederModule interface {
	Seed(entities []string, is_new_only bool) error
	GetSeedables() (entities []string, err error)
}
