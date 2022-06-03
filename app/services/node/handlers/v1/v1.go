// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/DMV-Petri-Dish/crypto/app/services/node/handlers/v1/private"
	"github.com/DMV-Petri-Dish/crypto/app/services/node/handlers/v1/public"
	"github.com/DMV-Petri-Dish/crypto/foundation/blockchain/state"
	"github.com/DMV-Petri-Dish/crypto/foundation/events"
	"github.com/DMV-Petri-Dish/crypto/foundation/nameservice"
	"github.com/DMV-Petri-Dish/crypto/foundation/web"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const version = "v1"

// config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *zap.SugaredLogger
	State *state.State
	NS    *nameservice.NameService
	Evts  *events.Events
}

// PublicRoutes binds all the version 1 public routes.
func PublicRoutes(app *web.App, cfg Config) {
	pbl := public.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
		WS:    websocket.Upgrader{},
		Evts:  cfg.Evts,
	}

	app.Handle(http.MethodGet, version, "/events", pbl.Events)
	app.Handle(http.MethodGet, version, "/genesis/list", pbl.Genesis)

	app.Handle(http.MethodPost, version, "/tx/submit", pbl.SubmitWalletTransaction)

}

// PrivateRoutes binds all the version 1 private routes
func PrivateRoutes(app *web.App, cfg Config) {
	prv := private.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	}

	app.Handle(http.MethodPost, version, "/node/peers", prv.SubmitPeer)
}
