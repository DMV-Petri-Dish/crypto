// package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

// A Handler is a type that handles a http request within our framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration data
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handles a set of routes for the application
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	// Create an OpenTelemetry HTTP Handler which wraps our router. This will start
	// the initial span and annotate it with information about the request/response
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if a client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/

	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// SignalShutdown is used to gracefully shutdown the app when an integrity issue is identified
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair to the server mux
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler
	handler = wrapMiddlware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddlware(a.mw, handler)

	// The function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {

		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Set the context with the required values to process the request
		v := Values{
			TraceID: uuid.New().String(),
			Now:     time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, key, &v)

		// Call the wrapped handler functions
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.ContextMux.Handle(method, finalPath, h)
}
