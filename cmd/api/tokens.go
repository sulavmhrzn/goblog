package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/sulavmhrzn/internal/data"
	"github.com/sulavmhrzn/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err.Error())
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePlaintextPassword(v, input.Password)
	if !v.IsValid() {
		app.failedValidationCheckErrorResponse(w, r, v.Error)
		return
	}

	user, err := app.models.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.invalidCredentialsErrorResponse(w, r)
			return
		default:
			app.internalServerErrorResponse(w, r, err.Error())
			return
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	if !match {
		app.invalidCredentialsErrorResponse(w, r)
		return
	}
	token, err := app.models.TokenModel.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	err = app.writeJSON(w, r, envelope{"authentication_token": token}, http.StatusCreated)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
	}

}
