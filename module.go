package ginfx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.prilax.in/fxmodules/ginfx.git/config"
	"gitlab.prilax.in/fxmodules/ginfx.git/middleware"
	"go.uber.org/fx"
)

type Route interface {
	Register(router *gin.Engine)
}

type Params struct {
	fx.In
	Config config.ServerConfig // Injected from configfx
	Routes []Route             `group:"routes"`
}

var Module = fx.Module("ginfx",
	fx.Provide(NewGinEngine),
	fx.Provide(middleware.NewValidator),
	fx.Invoke(func(lc fx.Lifecycle, r *gin.Engine, params Params) {
		for _, route := range params.Routes {
			route.Register(r)
		}
		fmt.Println()
		fmt.Printf("\nGin server is starting at port: %s\n", params.Config.GetPort())
		fmt.Printf("\nKeycloakRealmURL: %s\n", params.Config.GetKeycloakRealmURL())
		fmt.Println()
		server := &http.Server{
			Addr:    ":" + params.Config.GetPort(),
			Handler: r,
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {

				go server.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}),
)

func NewGinEngine(cfg config.ServerConfig, v *middleware.Validator) *gin.Engine {
	gin.SetMode(cfg.GetMode())
	r := gin.New()
	// Disable the redirects that break CORS
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.JWTMiddleware(v))
	return r
}

// AsRoute tags a handler as a member of the "routes" group
func AsRoute(f any) any {
	return fx.Annotate(f, fx.As(new(Route)), fx.ResultTags(`group:"routes"`))
}
