package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/data"
	"github.com/jessicatarra/greenlight/internal/validator"
	"net/http"
	"time"
)

type createAuthTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Create authentication token
// @Description Creates an authentication token for a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body createAuthTokenRequest true "Request body"
// @Success 201 {object} data.Token "Authentication token"
// @Router /tokens/authentication [post]
func (app *application) createAuthenticationTokenHandler(writer http.ResponseWriter, request *http.Request) {
	input := createAuthTokenRequest{}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(writer, request)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
