package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message string, status int) {
	data := envelope{"error": message}
	app.writeJSON(w, r, data, status)
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	msg := "internal server error"
	app.errorlog.Println(message)
	app.errorResponse(w, r, msg, http.StatusInternalServerError)
}
