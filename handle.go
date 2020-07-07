package aqua

import (
	"net/http"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables) and returns an error to ba handled by framework
type Handle func(http.ResponseWriter, *http.Request, Params) error
