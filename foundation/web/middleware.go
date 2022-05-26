package web

// Middleware is a function designed to run some code before and/or after another Handler.
// It is designed to remove boilerplate or other concerns not direclt to a given Handler.
type Middleware func(Handler) Handler

// wrapMiddlware creates a new handler by wrapping middleware around a final handler.
// The middlewares' Handlers will be executed by requests in the order they are provided.
func wrapMiddlware(mw []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one.
	// Replace the handler with the new wrapped handler.
	// Looping backwards ensures that the first middleware of the slice
	// is the first to executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
