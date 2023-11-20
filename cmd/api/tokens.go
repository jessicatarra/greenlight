package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/validator"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"strconv"
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

	database.ValidateEmail(v, input.Email)
	database.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
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

	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(user.ID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "greenlight.tarralva.com"
	claims.Audiences = []string{"greenlight.tarralva.com"}
	claims.Audiences = []string{"greenlight.tarralva.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
