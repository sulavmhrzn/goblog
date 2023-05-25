package main

import "net/http"

func (app *application) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID := app.contextGetUser(r)
	dashboard, err := app.models.UserModel.DashboardDetails(userID.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err.Error())
		return
	}
	app.writeJSON(w, r, envelope{"dashboard": dashboard}, http.StatusOK)
}
