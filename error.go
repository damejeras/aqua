package aqua

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
