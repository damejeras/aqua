package aqua

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Middleware func(next Handle) Handle

func chainMiddleware(handle Handle, mw ...Middleware) Handle {
	return func(w http.ResponseWriter, rq *http.Request, p httprouter.Params) error {
		if handle == nil {
			handle = func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
				return nil
			}
		}

		for i := range mw {
			handle = mw[len(mw)-1-i](handle)
		}

		return handle(w, rq, p)
	}
}
