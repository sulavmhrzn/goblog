package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sulavmhrzn/internal/data"
	"github.com/sulavmhrzn/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenErrorResponse(w, r)
			return
		}
		token := headerParts[1]

		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.IsValid() {
			app.invalidAuthenticationTokenErrorResponse(w, r)
			return
		}

		user, err := app.models.UserModel.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRows):
				app.invalidCredentialsErrorResponse(w, r)
			default:
				app.internalServerErrorResponse(w, r, err.Error())
			}
			return
		}
		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)

	})
}
