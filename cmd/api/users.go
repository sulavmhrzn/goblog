package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sulavmhrzn/goblog/internal/data"
	"github.com/sulavmhrzn/goblog/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err.Error())
		return
	}

	user := &data.User{
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}

	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsValid() {
		app.failedValidationCheckErrorResponse(w, r, v.Error)
		return
	}

	err = app.models.UserModel.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.badRequestErrorResponse(w, r, err.Error())
			return
		default:
			app.internalServerErrorResponse(w, r, err.Error())
			return
		}
	}
	token, err := app.models.TokenModel.New(user.ID, 24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	app.background(func() {
		body := fmt.Sprintf("Your activation token is: %s", token.Plaintext)
		err = app.mailer.Send(user.Email, "Welcome to Goblog!", body)
		if err != nil {
			app.internalServerErrorResponse(w, r, err.Error())
			return
		}
	})
	app.writeJSON(w, r, envelope{"data": user}, http.StatusOK)
}
