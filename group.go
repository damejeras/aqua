package aqua

type group struct {
	prefix string
	middleware []Handle
	parentRouter Router
}

func (r *group) NewGroup(prefix string) Router {
	return &group{prefix, make([]Handle, 0), r}
}

func (r *group) Use(handle Handle) {
	r.middleware = append(r.middleware, handle)
}

func (r *group) GET(path string, handle Handle, mw ...Handle) {
	r.parentRouter.GET(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r *group) HEAD(path string, handle Handle, mw ...Handle) {
	r.parentRouter.HEAD(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r *group) OPTIONS(path string, handle Handle, mw ...Handle) {
	r.parentRouter.OPTIONS(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r *group) POST(path string, handle Handle, mw ...Handle) {
	r.parentRouter.POST(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r *group) PUT(path string, handle Handle, mw ...Handle) {
	r.parentRouter.PUT(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r group) PATCH(path string, handle Handle, mw ...Handle) {
	r.parentRouter.PATCH(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}

func (r *group) DELETE(path string, handle Handle, mw ...Handle) {
	r.parentRouter.DELETE(r.prefix + path, chanHandles(handle, append(r.middleware, mw...)...))
}
