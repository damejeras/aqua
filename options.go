package aqua

import "net/http"

type Option func(s *rootRouter)

func DisableMethodNotAllowed() Option {
	return func(rr *rootRouter) {
		rr.HandleMethodNotAllowed = false
	}
}

func WithCustomMethodNotAllowedHandler(handler http.Handler) Option {
	return func(rr *rootRouter) {
		rr.HandleMethodNotAllowed = true
		rr.MethodNotAllowed = handler
	}
}

func WithCustomNotFoundHandler(handler http.Handler) Option {
	return func(rr *rootRouter) {
		rr.NotFound = handler
	}
}

func DisableErrorHandling() Option {
	return func(rr *rootRouter) {
		rr.config.ErrorHandler = func(next Handle) Handle {
			return func(w http.ResponseWriter, r *http.Request, p Params) error {
				return next(w, r, p)
			}
		}
	}
}

func WithCustomErrorHandler(handler Middleware) Option {
	return func(rr *rootRouter) {
		rr.config.ErrorHandler = handler
	}
}
