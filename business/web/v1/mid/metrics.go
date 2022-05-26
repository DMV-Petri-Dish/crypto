package mid

import (
	"context"
	"net/http"

	"github.com/DMV-Petri-Dish/crypto/business/web/metrics"
	"github.com/DMV-Petri-Dish/crypto/foundation/web"
)

// Metrics updates program counters
func Metrics() web.Middleware {

	// This is the actual middleware function to be executed
	m := func(hanlder web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Add the metrics into the context for metric gathering.
			ctx = metrics.Set(ctx)

			// Call the next handler.
			err := handler(ctx, w, r)

			// Handle updating the metrics that can be handled here.

			// Increment the request and goroutines counter.
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			// Increment if there is an error flowing through the request.
			if err != nil {
				metrics.AddErrors(ctx)
			}

			// Return the error so it can be handled futher up the chain
			return err
		}

		return h
	}

	return m
}
