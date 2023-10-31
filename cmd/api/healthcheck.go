package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	json, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	json = append(json, '\n')

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(json)
}
