package aqua

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	http.Handler

	Use(mw ...Middleware)
	NewGroup(prefix string) Router
	GET(path string, handle Handle, mw ...Middleware)
	HEAD(path string, handle Handle, mw ...Middleware)
	OPTIONS(path string, handle Handle, mw ...Middleware)
	POST(path string, handle Handle, mw ...Middleware)
	PUT(path string, handle Handle, mw ...Middleware)
	PATCH(path string, handle Handle, mw ...Middleware)
	DELETE(path string, handle Handle, mw ...Middleware)
}

func NewRouter(options ...Option) Router {
	router := httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound:               defaultNotFoundHandler,
		MethodNotAllowed:       defaultMethodNotAllowedHandler,
	}

	wrappedRouter := rootRouter{
		Router:     &router,
		middleware: make([]Middleware, 0),
		config: &Config{
			ErrorHandler: defaultErrorHandler,
		},
	}

	for _, option := range options {
		option(&wrappedRouter)
	}

	return &wrappedRouter
}

type rootRouter struct {
	*httprouter.Router

	middleware []Middleware

	config *Config
}

func (rr *rootRouter) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	rr.Router.ServeHTTP(w, rq)
}

func (rr *rootRouter) Use(mw ...Middleware) {
	rr.middleware = append(rr.middleware, mw...)
}

func (rr *rootRouter) NewGroup(prefix string) Router {
	return &group{prefix, make([]Middleware, 0), rr}
}

func (rr *rootRouter) GET(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodGet, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) HEAD(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodHead, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) OPTIONS(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodOptions, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) POST(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodPost, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) PUT(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodPut, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) PATCH(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodPatch, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) DELETE(path string, handle Handle, mw ...Middleware) {
	rr.handle(http.MethodDelete, path, chainMiddleware(handle, append(rr.middleware, mw...)...))
}

func (rr *rootRouter) handle(method, path string, handle Handle) {
	rr.Handle(method, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_ = rr.config.ErrorHandler(handle)(w, r, Params{p})
	})
}
