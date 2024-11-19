package modules

import (
	"fmt"

	"github.com/elmawardy/nutrix/common/customerrors"
	"github.com/elmawardy/nutrix/common/logger"
)

type AppManager struct {
	Modules map[string]IBaseModule
	Logger  logger.ILogger
}

func (manager *AppManager) Ignite() (err error) {

	for _, saved_module_builder := range saved_module_builders {
		manager.IgniteModule("core", manager.Logger, saved_module_builder)
	}

	return err
}

func (manager *AppManager) LoadModule(module IBaseModule, module_name string) *ModuleBuilder {

	builder := &ModuleBuilder{
		Logger: manager.Logger,
	}

	builder.module = module
	builder.module_name = module_name

	return builder
}

func (manager *AppManager) IgniteModule(name string, logger logger.ILogger, module_builder *ModuleBuilder) error {

	if manager.Modules == nil {
		manager.Modules = make(map[string]IBaseModule)
	}

	msg := fmt.Sprintf(`Running module : (%s) ...`, name)
	logger.Info(msg)

	if _, ok := manager.Modules[name]; ok {
		return customerrors.ErrModuleNameAlreadyExists
	}

	manager.Modules[name] = module_builder.module

	if module_builder.isRegisterHttpHandlers {
		if m, ok := module_builder.module.(IHttpModule); ok {
			m.RegisterHttpHandlers(module_builder.httpRouter, "/"+name)
		} else {
			logger.Error(customerrors.ErrTypeAssersionFailed.Error())
		}
	}

	if module_builder.isRegisterBackgroundWorkers {

		if m, ok := module_builder.module.(IBackgroundWorkerModule); ok {

			bw_svc := &background_worker_svc{
				Logger:  logger,
				Workers: m.RegisterBackgroundWorkers(),
			}

			go bw_svc.Start()
		}

	}

	logger.Info("Successfully registered module (" + name + ")")

	return nil
}

func (manager *AppManager) GetModules() (modules map[string]IBaseModule, err error) {

	if manager.Modules == nil {
		return nil, customerrors.ErrModuleNotRegistered
	}

	return manager.Modules, nil
}
