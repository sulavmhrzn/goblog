package main

import (
	"encoding/json"
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
