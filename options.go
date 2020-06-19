package aqua

import "net/http"

type Option func(s *rootRouter)

// WithMethodNotAllowedHandler takes Handle as an argument. It will be executed on requests to endpoints
// which are registered request method doesn't match method that endpoint accepts.
func WithMethodNotAllowedHandler(handle Handle) Option {
	return func(rr *rootRouter) {
		rr.HandleMethodNotAllowed = true
		rr.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = chainMiddleware(handle, rr.middleware...)(w, r, Params{})
		})
	}
}

// WithNotFoundHandler takes Handle as an argument. It will be executed on requests to non existing endpoints.
// It is also wrapped into Middleware functions registered for the router.
func WithNotFoundHandler(handle Handle) Option {
	return func(rr *rootRouter) {
		rr.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = chainMiddleware(handle, rr.middleware...)(w, r, Params{})
		})
	}
}
