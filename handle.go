package aqua

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables) and returns an error to ba handled by framework
type Handle func(http.ResponseWriter, *http.Request, httprouter.Params) error

func chanHandles(handle Handle, mw ...Handle) Handle {
	return func(w http.ResponseWriter, rq *http.Request, p httprouter.Params) error {
		var err error

		for i := range mw {
			err = mw[len(mw)-1-i](w, rq, p)
			if err != nil {
				return err
			}
		}

		return handle(w, rq, p)
	}
}
