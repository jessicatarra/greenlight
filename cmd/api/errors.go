package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": request.Method,
		"request_url":    request.URL.String(),
	})
}

func (app *application) rateLimitExceededResponse(writer http.ResponseWriter, request *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(writer, request, http.StatusTooManyRequests, message)
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

func (app *application) badRequestResponse(writer http.ResponseWriter, request *http.Request, err error) {
	app.errorResponse(writer, request, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(writer http.ResponseWriter, request *http.Request, errors map[string]string) {
	app.errorResponse(writer, request, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(writer http.ResponseWriter, request *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(writer, request, http.StatusConflict, message)
}

func (app *application) invalidCredentialsResponse(writer http.ResponseWriter, request *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(writer, request, http.StatusUnauthorized, message)
}
