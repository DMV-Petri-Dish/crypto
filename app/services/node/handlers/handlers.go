// Package handlers manages the different versions of the API.
package handlers

import (
	"context"
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/DMV-Petri-Dish/crypto/app/services/node/handlers/debug/checkgrp"
	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/state"
	"github.com/DMV-Petri-Dish/crypto/foundation/nameservice"
	"github.com/docker/docker/daemon/events"
	"go.uber.org/zap"
)

// MuxConfig contains all the mandatory systems required by handlers
type MuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	State    *state.State
	NS       *nameservice.NameService
	Evts     *events.Events
}

// PublicMux constructs a http.Handler with all application routes defined.
func PublicMux(cfg MuxConfig) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Cors("*"),
		mid.Panics(),
	)

	// Accept CORS 'OPTIONS' preflight requests if config has been provided.
	// Don't forget to apply the CORS middleware to the routes that need it.
	// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	}
	app.Handle(http.MethodOptions, "", "/*", h, mid.Cors("*"))

	// Load the v1 routes.
	v1.PublicRoutes(app, v1.Config{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
		Evts:  cfg.Evts,
	})

	return app
}

// PrivateMux constructs a http.Handler with all application routes defined.
func PrivateMux(cfg MuxConfig) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Cors("*"),
		mid.Panics(),
	)

	// Accept CORS 'OPTIONS' preflight requests if config has been provided.
	// Don't forget to apply the CORS middleware to the routes that need it.
	// Example Config: `conf:"default:https://MY_DOMAIN.COM"`
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	}
	app.Handle(http.MethodOptions, "", "/*", h, mid.Cors("*"))

	// Load the v1 routes.
	v1.PrivateRoutes(app, v1.Config{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	})

	return app
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux, bypassing the use of DefaultServerMux.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all standard library debug endpoints
	// from the net/http/pprof package
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then
// debug application routes for the service. This bypasses the use of
// DefaultServerMux. Using the DefaultServerMux would be a security risk
// because a dependency could inject a handler into our service with out our knowledge.
func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Construct the web.App which holds all routes as well as common Middleware.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}
