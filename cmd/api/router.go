package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/", app.createUserHandler)
	return router
}
