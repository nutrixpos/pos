// Package modules provides a builder for initializing and running modules.
//
// The builder is used to:
// - Set the logger and prompter for the module.
// - Register HTTP handlers for the module.
// - Register background workers for the module.
// - Start the module.
package modules

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/pos/common/logger"
	"github.com/nutrixpos/pos/common/userio"
)

// ModuleBuilder is a builder for initializing and running modules.
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

// RegisterHttpHandlers registers HTTP handlers for the module.
func (builder *ModuleBuilder) RegisterHttpHandlers(router *mux.Router) *ModuleBuilder {

	builder.isRegisterHttpHandlers = true
	builder.httpRouter = router

	return builder
}

// RegisterBackgroundWorkers registers background workers for the module.
func (builder *ModuleBuilder) RegisterBackgroundWorkers() *ModuleBuilder {

	builder.isRegisterBackgroundWorkers = true

	return builder
}

// Save saves the module builder for later use.
func (builder *ModuleBuilder) Save() {

	saved_module_builders[builder.module_name] = builder

}
