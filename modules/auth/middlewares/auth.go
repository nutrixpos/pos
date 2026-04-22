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

type IAuthService interface {
	AllowAuthenticated(next http.Handler) http.Handler
	AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler
}

type NoAuth struct {
	Config config.Config
}

func NewNoAuth(conf config.Config) *NoAuth {
	return &NoAuth{Config: conf}
}

func (na *NoAuth) AllowAuthenticated(next http.Handler) http.Handler {
	return next
}

func (na *NoAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {
	return next
}

type InternalAuth struct {
	JWTUtil *JWTUtil
	Config  config.Config
}

func NewInternalAuth(conf config.Config, jwtUtil *JWTUtil) *InternalAuth {
	return &InternalAuth{
		JWTUtil: jwtUtil,
		Config:  conf,
	}
}

func (ia *InternalAuth) AllowAuthenticated(next http.Handler) http.Handler {
	if !ia.Config.Auth.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no authorization header found", http.StatusForbidden)
			return
		}

		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		claims, err := ia.JWTUtil.ValidateToken(token)
		if err != nil {
			h := ia.Config.Zitadel
			if h.Enabled {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "invalid token", http.StatusForbidden)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "auth_ctx", claims))
		next.ServeHTTP(w, r)
	})
}

func (ia *InternalAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ia.Config.Auth.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "no authorization header found", http.StatusForbidden)
			return
		}

		token := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		}

		claims, err := ia.JWTUtil.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		for _, userRole := range claims.Roles {
			if userRole == "superuser" {
				r = r.WithContext(context.WithValue(r.Context(), "auth_ctx", claims))
				next.ServeHTTP(w, r)
				return
			}
		}

		for _, role := range roles {
			for _, userRole := range claims.Roles {
				if userRole == role {
					r = r.WithContext(context.WithValue(r.Context(), "auth_ctx", claims))
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		http.Error(w, "insufficient permissions", http.StatusForbidden)
	})
}

type ZitadelAuth struct {
	Domain string
	Key    string
	AuthZ  *authorization.Authorizer[*oauth.IntrospectionContext]
	Config config.Config
}

func NewZitadelAuth(conf config.Config) *ZitadelAuth {
	if !conf.Zitadel.Enabled {
		return &ZitadelAuth{}
	}

	ctx := context.Background()

	za := ZitadelAuth{
		Domain: conf.Zitadel.Domain,
		Key:    conf.Zitadel.KeyPath,
		Config: conf,
	}

	portStr := strconv.Itoa(conf.Zitadel.Port)

	authZ, err := authorization.New(ctx, zitadel.New(za.Domain, zitadel.WithInsecure(portStr)), oauth.DefaultAuthorization(za.Key))

	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	za.AuthZ = authZ

	return &za
}

func (za *ZitadelAuth) AllowAuthenticated(next http.Handler) http.Handler {
	if !za.Config.Zitadel.Enabled {
		return next
	}

	mw := middleware.New(za.AuthZ)
	handler := mw.RequireAuthorization()
	return handler(next)
}

func (za *ZitadelAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {
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

			authCtx, err := za.AuthZ.CheckAuthorization(r.Context(), reqToken, authorization.WithRole(role))
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), "auth_ctx", authCtx.IntrospectionResponse))
				authorized = true
				next.ServeHTTP(w, r)
				break
			}
		}

		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}