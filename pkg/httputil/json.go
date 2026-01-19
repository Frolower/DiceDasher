package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

func UnpackJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if dec.More() {
		return errors.New("invalid json: multiple objects")
	}

	return nil
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
