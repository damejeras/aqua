package aqua

type Error struct {
	Code    int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Cause   error                  `json:"-"`
}

func (err Error) Error() string {
	if err.Cause == nil {
		return err.Message
	}

	return err.Message + ": " + err.Cause.Error()
}

type ErrorLogger func(error)
