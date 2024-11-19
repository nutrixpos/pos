package auth

import (
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
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
