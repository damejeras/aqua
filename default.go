package aqua

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func defaultErrorLogger(err error) {
	log.Printf("ERROR: %v", err)
}

func defaultErrorHandler(next Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p Params) error {
		err := next(w, r, p)
		if err != nil {
			defaultErrorLogger(err)
			w.Header().Set("Content-Type", "application/json")
			if aquaErr, ok := err.(Error); ok {
				w.WriteHeader(aquaErr.Code)
				if _, err = w.Write([]byte(`{"status": ` +
					strconv.FormatInt(int64(aquaErr.Code), 10) +
					`, "message": "` +
					aquaErr.Message + `"}`)); err != nil {
					defaultErrorLogger(err)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				if _, err = w.Write([]byte(`{"status": 500, "message": "internal server error"}`)); err != nil {
					defaultErrorLogger(err)
				}
			}
		}

		return nil
	}
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_ = defaultErrorHandler(func(w http.ResponseWriter, r *http.Request, p Params) error {
		return Error{
			Code:    http.StatusNotFound,
			Message: "not found",
			Cause:   fmt.Errorf("%s", r.URL.Path),
		}
	})(w, r, Params{})
})

var defaultMethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_ = defaultErrorHandler(func(w http.ResponseWriter, r *http.Request, p Params) error {
		return Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
			Cause:   fmt.Errorf("%s %s", r.Method, r.URL.Path),
		}
	})(w, r, Params{})
})
