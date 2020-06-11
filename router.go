package aqua

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router interface {
	http.Handler

	Use(handle Handle)
	NewGroup(prefix string) Router
	GET(path string, handle Handle, mw ...Handle)
	HEAD(path string, handle Handle, mw ...Handle)
	OPTIONS(path string, handle Handle, mw ...Handle)
	POST(path string, handle Handle, mw ...Handle)
	PUT(path string, handle Handle, mw ...Handle)
	PATCH(path string, handle Handle, mw ...Handle)
	DELETE(path string, handle Handle, mw ...Handle)
}

func NewRouter() Router {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(encodeMessage("Resource not found"))
	})

	return &rootRouter{r, make([]Handle, 0)}
}

type rootRouter struct {
	*httprouter.Router

	middleware []Handle
}

func (r *rootRouter) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	r.Router.ServeHTTP(w, rq)
}

func (r *rootRouter) Use(handle Handle) {
	r.middleware = append(r.middleware, handle)
}

func (r *rootRouter) NewGroup(prefix string) Router {
	return &group{prefix, make([]Handle, 0), r}
}

func (r *rootRouter) GET(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodGet, path, chanHandles(handle, mw...))
}

func (r *rootRouter) HEAD(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodHead, path, chanHandles(handle, mw...))
}

func (r *rootRouter) OPTIONS(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodOptions, path, chanHandles(handle, mw...))
}

func (r *rootRouter) POST(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodPost, path, chanHandles(handle, mw...))
}

func (r *rootRouter) PUT(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodPut, path, chanHandles(handle, mw...))
}

func (r *rootRouter) PATCH(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodPatch, path, chanHandles(handle, mw...))
}

func (r *rootRouter) DELETE(path string, handle Handle, mw ...Handle) {
	r.handle(http.MethodDelete, path, chanHandles(handle, mw...))
}

func (r *rootRouter) handle(method, path string, handle Handle) {
	r.Router.Handle(method, path, func(w http.ResponseWriter, rq *http.Request, p httprouter.Params) {
		err := chanHandles(handle, r.middleware...)(w, rq, p)
		if err == nil {
			return
		}

		log.Printf("ERROR: %v", err)

		clientError, ok := err.(ClientError)
		if !ok {
			// if it is not ClientError, assume that it is ServerError.
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(encodeMessage("Internal server error"))
			return
		}

		w.WriteHeader(clientError.ResponseStatus())
		_, _ = w.Write(clientError.ResponseBody())
	})
}
