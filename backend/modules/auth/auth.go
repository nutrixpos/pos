package auth

import (
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
	"github.com/elmawardy/nutrix/modules"
	"github.com/gorilla/mux"
)

type IHttpAuth interface {
	AllowRoles(next http.Handler, roles ...string) http.Handler
	AllowAuthenticated(next http.Handler) http.Handler
}

func NewBuilder(config config.Config, settings config.Settings) *AuthModuleBuilder {
	mb := new(AuthModuleBuilder)
	mb.Config = config
	mb.Settings = settings

	return mb
}

type Auth struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

type AuthModuleBuilder struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

func (amb *AuthModuleBuilder) SetLogger(logger logger.ILogger) modules.IModuleBuilder {
	amb.Logger = logger
	return amb
}

func (amb *AuthModuleBuilder) SetPrompter(prompter userio.Prompter) modules.IModuleBuilder {
	amb.Prompter = prompter
	return amb
}

func (amb *AuthModuleBuilder) RegisterHttpHandlers(router *mux.Router, prefix string) modules.IModuleBuilder {

	return amb
}

func (amb *AuthModuleBuilder) Build() modules.BaseModule {
	return &Auth{
		Logger:   amb.Logger,
		Config:   amb.Config,
		Prompter: amb.Prompter,
		Settings: amb.Settings,
	}

}
