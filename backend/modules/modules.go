package modules

import (
	"fmt"

	"github.com/elmawardy/nutrix/common/customerrors"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/gorilla/mux"
)

type ModulesManager struct {
	Modules map[string]IModule
}

type IHttpModuleBuilder interface {
	RegisterHttpHandlers(*mux.Router) IModuleBuilder
}

type IModuleBuilder interface {
	SetLogger(logger.ILogger) IModuleBuilder
	Build() IModule
}

type IModule interface {
}

func (manager *ModulesManager) RegisterModule(name string, logger logger.ILogger, module_builder IModuleBuilder, r ...*mux.Router) error {

	if manager.Modules == nil {
		manager.Modules = make(map[string]IModule)
	}

	msg := fmt.Sprintf("Registering module %s: ", name)
	logger.Info(msg)

	if _, ok := manager.Modules[name]; ok {
		return customerrors.ModuleNameAlreadyExists{}
	}

	new_module := module_builder.SetLogger(logger)

	if len(r) > 0 {
		if m, ok := module_builder.(IHttpModuleBuilder); ok {
			m.RegisterHttpHandlers(r[0])
		} else {
			logger.Error(customerrors.TypeAssersionFailed{}.Error())
		}
	}

	manager.Modules[name] = new_module.Build()

	return nil
}
