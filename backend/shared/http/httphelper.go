package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// verifies correctness of JSON file and returns a Go object as well as corresponding error
func ReadJSON[T any](r *http.Request) (T, error) {
	var obj T

	maxBytes := 1_048_576 // 1MB max JSON size
	r.Body = http.MaxBytesReader(nil, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&obj); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return obj, fmt.Errorf("malformed JSON at position %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return obj, fmt.Errorf("malformed JSON")
		case errors.As(err, &unmarshalTypeError):
			return obj, fmt.Errorf("invalid value for field %q", unmarshalTypeError.Field)
		case errors.Is(err, io.EOF):
			return obj, fmt.Errorf("body must not be empty")
		default:
			return obj, err
		}
	}

	if decoder.More() {
		return obj, fmt.Errorf("body must only contain a single JSON object or array")
	}

	return obj, nil
}

// Converts Go object into a JSON file ready to be sent further
func WriteJSON[T any](status int, data T) ([]byte, int, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, 0, err
	}

	return jsonBytes, status, nil
}
