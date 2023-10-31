package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"status":      "available",
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(writer, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
