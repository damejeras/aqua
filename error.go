package aqua

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type ClientError interface {
	Error() string
	ResponseBody() []byte
	ResponseStatus() int
}

func Err(err error, status int, message string) error {
	return &Error{
		Cause:      err,
		Message:    message,
		StatusCode: status,
	}
}

type Error struct {
	Cause      error
	Message    string
	StatusCode int
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + " : " + e.Cause.Error()
}

func (e *Error) ResponseBody() []byte {
	return encodeMessage(e.Message)
}

func (e *Error) ResponseStatus() int {
	return e.StatusCode
}

func encodeMessage(message string) []byte {
	return []byte(`{"message": "` + message + `"}`)
}

func defaultErrorHandler(next Handle) Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
		err := next(w, r, p)
		if err != nil {
			log.Printf("ERROR: %v", err)

			clientError, ok := err.(ClientError)
			if !ok {
				// if it is not ClientError, assume that it is ServerError.
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(internalErrorMessage)
				return nil
			}

			w.WriteHeader(clientError.ResponseStatus())
			_, _ = w.Write(clientError.ResponseBody())
		}

		return nil
	}
}
