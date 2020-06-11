package aqua

import "net/http"

type group struct {
	prefix       string
	middleware   []Middleware
	parentRouter Router
}

func (g *group) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	g.parentRouter.ServeHTTP(w, rq)
}

func (g *group) NewGroup(prefix string) Router {
	return &group{prefix, make([]Middleware, 0), g}
}

func (g *group) Use(mw ...Middleware) {
	g.middleware = append(g.middleware, mw...)
}

func (g *group) GET(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.GET(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g *group) HEAD(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.HEAD(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g *group) OPTIONS(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.OPTIONS(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g *group) POST(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.POST(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g *group) PUT(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.PUT(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g group) PATCH(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.PATCH(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}

func (g *group) DELETE(path string, handle Handle, mw ...Middleware) {
	g.parentRouter.DELETE(g.prefix+path, chainMiddleware(handle, append(g.middleware, mw...)...))
}
