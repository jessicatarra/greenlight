package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(writer http.ResponseWriter, request *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(writer, status, env, nil)
	if err != nil {
		app.logError(request, err)
		writer.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(writer, request, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(writer http.ResponseWriter, request *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(writer, request, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(writer http.ResponseWriter, request *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", request.Method)
	app.errorResponse(writer, request, http.StatusMethodNotAllowed, message)
}
