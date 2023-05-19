package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) router() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/", app.createUserHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/activate", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/activated", app.activateUserHandler)
	return app.authenticate(router)
}
