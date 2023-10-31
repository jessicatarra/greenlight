package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJSON(writer, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
