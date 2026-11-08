package httputil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"net/http"
)

func UnpackJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		// no body at all, treat like empty
		if _, ok := dst.(*jsontext.Value); ok {
			return nil
		}
		return errors.New("empty body")
	}
	defer r.Body.Close()

	dec := jsontext.NewDecoder(r.Body)

	// Decode first JSON value
	if err := json.UnmarshalDecode(dec, dst, json.RejectUnknownMembers(true)); err != nil {
		// Empty body -> EOF
		if errors.Is(err, io.EOF) {
			// For jsontext.Value, empty is acceptable; leave it nil
			if _, ok := dst.(*jsontext.Value); ok {
				return nil
			}
			return errors.New("empty body")
		}
		return err
	}

	// Ensure there's no trailing junk / extra JSON values
	if err := json.UnmarshalDecode(dec, &struct{}{}); err != nil {
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

	enc := jsontext.NewEncoder(w, jsontext.EscapeForHTML(true))

	return json.MarshalEncode(enc, src)
}
