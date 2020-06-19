package aqua

import (
	"net/http"
)

type Middleware func(next Handle) Handle

func chainMiddleware(handle Handle, mw ...Middleware) Handle {
	if handle == nil {
		handle = func(w http.ResponseWriter, r *http.Request, p Params) error {
			return nil
		}
	}

	for i := range mw {
		handle = mw[len(mw)-1-i](handle)
	}

	return handle
}
