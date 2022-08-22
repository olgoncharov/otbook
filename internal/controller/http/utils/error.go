package utils

import (
	"errors"
	std_http "net/http"

	jsoniter "github.com/json-iterator/go"
)

var (
	ErrInternal = errors.New("internal error")
)

type HTTPError struct {
	Error string `json:"error"`
}

func WriteJSONError(w std_http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	encoder := jsoniter.NewEncoder(w)
	encoder.Encode(HTTPError{
		Error: err,
	})
}
