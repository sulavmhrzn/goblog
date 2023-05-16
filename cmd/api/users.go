package main

import (
	"net/http"

	"github.com/sulavmhrzn/internal/data"
	"github.com/sulavmhrzn/internal/validator"
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
	app.writeJSON(w, r, envelope{"data": input}, http.StatusOK)
}
