package aqua

import (
	"fmt"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidParamType = fmt.Errorf("unexpected param type")
	ErrMissingParam     = fmt.Errorf("missing param")
)

// Params wraps httprouter.Params and provides additional param type conversion methods
type Params struct {
	httprouter.Params
}

func (p Params) Int(name string) (int, error) {
	v := p.ByName(name)
	if v == "" {
		return 0, ErrMissingParam
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, ErrInvalidParamType
	}

	return i, nil
}

func (p Params) String(name string) (string, error) {
	v := p.ByName(name)
	if v == "" {
		return "", ErrMissingParam
	}

	return v, nil
}

func (p Params) Float64(name string) (float64, error) {
	v := p.ByName(name)
	if v == "" {
		return 0, ErrMissingParam
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, ErrInvalidParamType
	}

	return f, nil
}
