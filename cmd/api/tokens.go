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
		return
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

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	v := validator.New()

	data.ValidatePlaintextPassword(v, input.Password)
	data.ValidateEmail(v, input.Email)
	if !v.IsValid() {
		app.failedValidationCheckErrorResponse(w, r, v.Error)
		return
	}

	user, err := app.models.UserModel.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.invalidCredentialsErrorResponse(w, r)
		default:
			app.internalServerErrorResponse(w, r, err.Error())
		}
		return
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
	token, err := app.models.TokenModel.New(user.ID, 24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	err = app.writeJSON(w, r, envelope{"activation_token": token}, http.StatusCreated)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err.Error())
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.IsValid() {
		app.failedValidationCheckErrorResponse(w, r, v.Error)
		return
	}

	user, err := app.models.UserModel.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			v.AddErrorMessage("token", "invalid or expired activation token")
			app.failedValidationCheckErrorResponse(w, r, v.Error)
			return
		default:
			app.internalServerErrorResponse(w, r, err.Error())
			return
		}
	}
	user.Activated = true

	err = app.models.UserModel.Update(user)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	err = app.models.TokenModel.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	err = app.writeJSON(w, r, envelope{"user": user}, http.StatusOK)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
	}
}
