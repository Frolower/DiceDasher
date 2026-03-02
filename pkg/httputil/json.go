package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func UnpackJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		// no body at all, treat like empty
		if _, ok := dst.(*json.RawMessage); ok {
			return nil
		}
		return errors.New("empty body")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode first JSON value
	if err := dec.Decode(dst); err != nil {
		// Empty body -> EOF
		if errors.Is(err, io.EOF) {
			// For RawMessage, empty is acceptable; leave it nil
			if _, ok := dst.(*json.RawMessage); ok {
				return nil
			}
			return errors.New("empty body")
		}
		return err
	}

	// Ensure there's no trailing junk / extra JSON values
	if err := dec.Decode(&struct{}{}); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.New("invalid json: trailing data")
	}

	return errors.New("invalid json: multiple values")
}

func PackJSON(w http.ResponseWriter, status int, src any) error {
	if src == nil {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	return enc.Encode(src)
}
