package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/elmawardy/nutrix/backend/common/config"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func NewZitadelAuth(conf config.Config) ZitadelAuth {

	ctx := context.Background()

	za := ZitadelAuth{
		Domain: "localhost",
		Key:    "zitadel-key.json",
	}

	authZ, err := authorization.New(ctx, zitadel.New(za.Domain, zitadel.WithInsecure("2020")), oauth.DefaultAuthorization(za.Key))
	/******  6999d1a5-4501-4af6-9052-14fbce64d1ab  *******/
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	za.AuthZ = authZ

	return za
}

type ZitadelAuth struct {
	Domain string // Zitadel instance domain
	Key    string // path to key.json
	AuthZ  *authorization.Authorizer[*oauth.IntrospectionContext]
}

func (za *ZitadelAuth) AllowAuthenticated(next http.Handler) http.Handler {

	mw := middleware.New(za.AuthZ)

	handler := mw.RequireAuthorization()

	return handler(next)
}

func (za *ZitadelAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {

	// Initialize the HTTP middleware by providing the authorization
	// mw := middleware.New(za.AuthZ)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
