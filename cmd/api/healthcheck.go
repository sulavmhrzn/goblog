package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]interface{}{"version": "1.0.0", "environment": "development"}
	err := app.writeJSON(w, r, envelope{"data": msg}, http.StatusOK)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
	}
}
