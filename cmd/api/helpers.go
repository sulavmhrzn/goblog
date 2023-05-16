package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type envelope map[string]interface{}

// writeJSON is a helper method that Marshals data, sets Content-Type to application/json and writes json to response
func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, data envelope, status int) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(js))
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		default:
			return err
		}
	}
	return nil
}
