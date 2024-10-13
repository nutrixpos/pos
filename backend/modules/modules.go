package modules

import (
	"fmt"

	"github.com/elmawardy/nutrix/common/customerrors"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/gorilla/mux"
)

type ModulesManager struct {
	Modules map[string]BaseModule
}

type IHttpModuleBuilder interface {
	RegisterHttpHandlers(*mux.Router) IModuleBuilder
}

type IModuleBuilder interface {
	SetLogger(logger.ILogger) IModuleBuilder
	SetPrompter(userio.Prompter) IModuleBuilder
	Build() BaseModule
}

type BaseModule interface {
}

type SeederModule interface {
	Seed(entities []string, is_new_only bool) error
	GetSeedables() (entities []string, err error)
}

func (manager *ModulesManager) RegisterModule(name string, logger logger.ILogger, module_builder IModuleBuilder, prompter userio.Prompter, r ...*mux.Router) error {

	if manager.Modules == nil {
		manager.Modules = make(map[string]BaseModule)
	}

	msg := fmt.Sprintf("Registering module : %s", name)
	logger.Info(msg)

	if _, ok := manager.Modules[name]; ok {
		return customerrors.ErrModuleNameAlreadyExists
	}

	new_module_builder := module_builder.SetLogger(logger).SetPrompter(prompter)

	if len(r) > 0 {
		if m, ok := module_builder.(IHttpModuleBuilder); ok {
			m.RegisterHttpHandlers(r[0])
		} else {
			logger.Error(customerrors.ErrTypeAssersionFailed.Error())
		}
	}

	manager.Modules[name] = new_module_builder.Build()

	return nil
}

func (manager *ModulesManager) GetModules() (modules map[string]BaseModule, err error) {

	if manager.Modules == nil {
		return nil, customerrors.ErrModuleNotRegistered
	}

	return manager.Modules, nil
}
