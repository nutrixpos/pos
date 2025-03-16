// Package middlewares provides a set of middleware functions used to check
// Zitadel access token for auth and roles.
package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/nutrixpos/pos/common/config"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

// NewZitadelAuth creates a new ZitadelAuth object with the given configuration.
// It sets up the Zitadel SDK with the given domain and key path.
func NewZitadelAuth(conf config.Config) ZitadelAuth {

	if !conf.Zitadel.Enabled {
		return ZitadelAuth{}
	}

	ctx := context.Background()

	za := ZitadelAuth{
		Domain: conf.Zitadel.Domain,
		Key:    conf.Zitadel.KeyPath,
	}

	portStr := strconv.Itoa(conf.Zitadel.Port)

	authZ, err := authorization.New(ctx, zitadel.New(za.Domain, zitadel.WithInsecure(portStr)), oauth.DefaultAuthorization(za.Key))

	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	za.AuthZ = authZ

	return za
}

// ZitadelAuth holds the configuration for Zitadel and the Authorizer
type ZitadelAuth struct {
	Domain string // Zitadel instance domain
	Key    string // path to key.json
	AuthZ  *authorization.Authorizer[*oauth.IntrospectionContext]
	Config config.Config
}

// AllowAuthenticated middleware checks if the given request has a valid acess token.
func (za *ZitadelAuth) AllowAuthenticated(next http.Handler) http.Handler {

	if !za.Config.Zitadel.Enabled {
		return next
	}

	mw := middleware.New(za.AuthZ)

	handler := mw.RequireAuthorization()

	return handler(next)
}

// AllowAnyOfRoles middleware checks if the given request has a valid access token
// and if the user has any of the given roles.
func (za *ZitadelAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {

	// Initialize the HTTP middleware by providing the authorization
	// mw := middleware.New(za.AuthZ)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !za.Config.Zitadel.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		authorized := false

		for _, role := range roles {

			reqToken := r.Header.Get("Authorization")
			if reqToken == "" {
				http.Error(w, "no authorization header found", http.StatusForbidden)
				return
			}

			_, err := za.AuthZ.CheckAuthorization(r.Context(), reqToken, authorization.WithRole(role))

			if err == nil {
				authorized = true
				next.ServeHTTP(w, r)
				break
			}
		}

		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
		}

	})

	// checkpoints := []authorization.CheckOption{}

	// for _, role := range roles {
	// 	checkpoints = append(checkpoints, authorization.WithRole(role))
	// }

	// handler := mw.RequireAuthorization(checkpoints...)

	// return handler(next)
}
