package aqua

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
