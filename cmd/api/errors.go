package main

import (
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message interface{}, status int) {
	data := envelope{"error": message}
	app.writeJSON(w, r, data, status)
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	msg := "internal server error"
	app.errorlog.Println(message)
	app.errorResponse(w, r, msg, http.StatusInternalServerError)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	msg := message
	app.errorResponse(w, r, msg, http.StatusBadRequest)
}

func (app *application) failedValidationCheckErrorResponse(w http.ResponseWriter, r *http.Request, message interface{}) {
	app.errorResponse(w, r, message, http.StatusUnprocessableEntity)
}

func (app *application) invalidCredentialsErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid credentials"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}

func (app *application) invalidAuthenticationTokenErrorResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, message, http.StatusUnauthorized)
}
