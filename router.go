package aqua

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	notFoundMessage      = encodeMessage("Resource not found")
	internalErrorMessage = encodeMessage("Internal server error")
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

func NewRouter() Router {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(notFoundMessage)
	})

	return &rootRouter{r, make([]Middleware, 0)}
}

type rootRouter struct {
	*httprouter.Router

	middleware []Middleware
}

func (r *rootRouter) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	r.Router.ServeHTTP(w, rq)
}

func (r *rootRouter) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

func (r *rootRouter) NewGroup(prefix string) Router {
	return &group{prefix, make([]Middleware, 0), r}
}

func (r *rootRouter) GET(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodGet, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) HEAD(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodHead, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) OPTIONS(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodOptions, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) POST(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodPost, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) PUT(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodPut, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) PATCH(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodPatch, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) DELETE(path string, handle Handle, mw ...Middleware) {
	r.handle(http.MethodDelete, path, chainMiddleware(handle, mw...))
}

func (r *rootRouter) handle(method, path string, handle Handle) {
	r.Router.Handle(method, path, func(w http.ResponseWriter, rq *http.Request, p httprouter.Params) {
		err := chainMiddleware(handle, r.middleware...)(w, rq, p)
		if err == nil {
			return
		}

		log.Printf("ERROR: %v", err)

		clientError, ok := err.(ClientError)
		if !ok {
			// if it is not ClientError, assume that it is ServerError.
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(internalErrorMessage)
			return
		}

		w.WriteHeader(clientError.ResponseStatus())
		_, _ = w.Write(clientError.ResponseBody())
	})
}
