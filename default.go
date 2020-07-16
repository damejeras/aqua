package aqua

import (
	"bytes"
	"io"
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
			w.Header().Set("Content-Type", "application/json")
			defaultErrorLogger(err)

			if aquaErr, ok := err.(Error); ok {
				w.WriteHeader(aquaErr.Code)
				if _, err = io.Copy(w, bytes.NewReader([]byte(`{"status": `+
					strconv.FormatInt(int64(aquaErr.Code), 10)+
					`, "message": "`+
					aquaErr.Message+`"}`))); err != nil {
					defaultErrorLogger(err)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				if _, err = io.Copy(
					w,
					bytes.NewReader([]byte(`{"status": 500, "message": "internal server error"}`)),
				); err != nil {
					defaultErrorLogger(err)
				}
			}
		}

		return nil
	}
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defaultErrorHandler(func(w http.ResponseWriter, r *http.Request, p Params) error {
		return &Error{
			Code:    http.StatusNotFound,
			Message: "not found",
			Cause:   nil,
		}
	})
})

var defaultMethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defaultErrorHandler(func(w http.ResponseWriter, r *http.Request, p Params) error {
		return &Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
			Cause:   nil,
		}
	})
})
