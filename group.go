package aqua

import "net/http"

type group struct {
	prefix       string
	middleware   []Handle
	parentRouter Router
}

func (g *group) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	g.parentRouter.ServeHTTP(w, rq)
}

func (g *group) NewGroup(prefix string) Router {
	return &group{prefix, make([]Handle, 0), g}
}

func (g *group) Use(handle Handle) {
	g.middleware = append(g.middleware, handle)
}

func (g *group) GET(path string, handle Handle, mw ...Handle) {
	g.parentRouter.GET(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g *group) HEAD(path string, handle Handle, mw ...Handle) {
	g.parentRouter.HEAD(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g *group) OPTIONS(path string, handle Handle, mw ...Handle) {
	g.parentRouter.OPTIONS(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g *group) POST(path string, handle Handle, mw ...Handle) {
	g.parentRouter.POST(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g *group) PUT(path string, handle Handle, mw ...Handle) {
	g.parentRouter.PUT(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g group) PATCH(path string, handle Handle, mw ...Handle) {
	g.parentRouter.PATCH(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}

func (g *group) DELETE(path string, handle Handle, mw ...Handle) {
	g.parentRouter.DELETE(g.prefix+path, chanHandles(handle, append(g.middleware, mw...)...))
}
