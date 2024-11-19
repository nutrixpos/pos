package modules

import (
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/gorilla/mux"
)

type ModuleBuilder struct {
	Logger                      logger.ILogger
	Prompter                    userio.Prompter
	module                      IBaseModule
	module_name                 string
	workers                     []Worker
	isRegisterHttpHandlers      bool
	isRegisterBackgroundWorkers bool
	httpRouter                  *mux.Router
}

func (builder *ModuleBuilder) RegisterHttpHandlers(router *mux.Router) *ModuleBuilder {

	builder.isRegisterHttpHandlers = true
	builder.httpRouter = router

	return builder
}

func (builder *ModuleBuilder) RegisterBackgroundWorkers() *ModuleBuilder {

	builder.isRegisterBackgroundWorkers = true

	return builder
}

func (builder *ModuleBuilder) Save() {

	saved_module_builders[builder.module_name] = builder

}
