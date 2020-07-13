package aqua

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

type Error struct {
	Code    int
	Message string
	Cause   error
}

func (err Error) Error() string {
	if err.Cause == nil {
		return err.Message
	}

	return err.Message + ": " + err.Cause.Error()
}

type ErrorLogger func(error)

func errorHandlingMiddleware(logger ErrorLogger) Middleware {
	return func(next Handle) Handle {
		if logger == nil {
			logger = func(err error) {}
		}

		return func(w http.ResponseWriter, r *http.Request, p Params) error {
			err := next(w, r, p)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				logger(err)

				if aquaErr, ok := err.(Error); ok {
					w.WriteHeader(aquaErr.Code)
					if _, err = io.Copy(w, bytes.NewReader([]byte(`{"status": `+
						strconv.FormatInt(int64(aquaErr.Code), 10)+
						`, "message": "`+
						aquaErr.Message+`"}`))); err != nil {
						logger(err)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					if _, err = io.Copy(
						w,
						bytes.NewReader([]byte(`{"status": 500, "message": "internal error"}`)),
					); err != nil {
						logger(err)
					}
				}
			}

			return nil
		}
	}
}
