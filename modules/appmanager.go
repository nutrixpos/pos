// Package modules is the package for managing the modules of the application.
//
// Modules are the building blocks of the application, and they can be used to
// extend the application with new features. The package provides functions and
// interfaces for managing the modules of the application.
package modules

import (
	"fmt"

	"github.com/nutrixpos/pos/common/customerrors"
	"github.com/nutrixpos/pos/common/logger"
)

// AppManager is a struct that manages the modules of the application.
//
// It contains a map of strings to IBaseModule, where the key is the name of the module
// and the value is the module itself. It also contains a logger.ILogger, which is
// used to log messages.
type AppManager struct {
	// Modules is a map of strings to IBaseModule, where the key is the name of the module
	// and the value is the module itself.
	Modules map[string]IBaseModule
	// Logger is a logger.ILogger, which is used to log messages.
	Logger logger.ILogger
}

// Run starts all saved module builders by igniting each module.
func (manager *AppManager) Run() (err error) {

	for _, saved_module_builder := range saved_module_builders {
		manager.RunModule(saved_module_builder.module_name, manager.Logger, saved_module_builder)
	}

	return err
}

// LoadModule initializes a ModuleBuilder with the specified module and name.
func (manager *AppManager) LoadModule(module IBaseModule, module_name string) *ModuleBuilder {

	builder := &ModuleBuilder{
		Logger: manager.Logger,
	}

	builder.module = module
	builder.module_name = module_name

	return builder
}

// RunModule initializes and registers a module from the module builder.
func (manager *AppManager) RunModule(name string, logger logger.ILogger, module_builder *ModuleBuilder) error {

	if manager.Modules == nil {
		manager.Modules = make(map[string]IBaseModule)
	}

	msg := fmt.Sprintf(`Starting module : (%s) ...`, name)
	logger.Info(msg)

	if _, ok := manager.Modules[name]; ok {
		return customerrors.ErrModuleNameAlreadyExists
	}

	manager.Modules[name] = module_builder.module

	err := module_builder.module.OnStart()()
	if err != nil {
		panic(err)
	}

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

	logger.Info("Started module (" + name + ")")

	return nil
}

// GetModules retrieves all registered modules.
func (manager *AppManager) GetModules() (modules map[string]IBaseModule, err error) {

	if manager.Modules == nil {
		return nil, customerrors.ErrModuleNotRegistered
	}

	return manager.Modules, nil
}
