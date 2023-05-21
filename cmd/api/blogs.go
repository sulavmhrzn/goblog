package main

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
	"github.com/sulavmhrzn/internal/data"
	"github.com/sulavmhrzn/internal/validator"
)

func (app *application) createBlogHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestErrorResponse(w, r, err.Error())
		return
	}

	v := validator.New()

	user := app.contextGetUser(r)
	blog := data.Blog{
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UserID:    user.ID,
		Slug:      slug.Make(input.Title),
	}

	if data.ValidateBlog(v, blog); !v.IsValid() {
		app.failedValidationCheckErrorResponse(w, r, v.Error)
		return
	}

	err = app.models.BlogModel.Insert(blog)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}

	app.writeJSON(w, r, envelope{"blog": blog}, http.StatusCreated)
}

func (app *application) listBlogsHandler(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.models.BlogModel.List()
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	app.writeJSON(w, r, envelope{"blogs": blogs}, http.StatusOK)
}