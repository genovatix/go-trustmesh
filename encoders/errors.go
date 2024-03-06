package encoders

import "errors"

var (
	ERROR_JSON_ENCODE = errors.New("there was an error encoding to JSON")
	ERROR_JSON_DECODE = errors.New("there was an error decoding from JSON")
)
