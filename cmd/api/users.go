package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/validator"
	"net/http"
	"time"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Activate User
// @Description Activates a user account using a token that was previously sent when successfully register a new user
// @Tags Users
// @Accept json
// @Produce  json
// @Param token query string true "Token for user activation"
// @Success 200 {object} data.User
// @Router /users/activated [put]
func (app *application) activateUserHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		TokenPlaintext string
	}

	qs := request.URL.Query()

	input.TokenPlaintext = app.readString(qs, "token", "")

	v := validator.New()

	if database.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(database.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrEditConflict):
			app.editConflictResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body createUserRequest true "User registration data"
// @Success 201 {object} data.User
// @Router /users [post]
func (app *application) registerUserHandler(writer http.ResponseWriter, request *http.Request) {
	input := createUserRequest{}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	user := &database.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	v := validator.New()

	if database.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, database.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		print(token.Plaintext)

		err = app.mailer.Send(user.Email, "user_welcome.gohtml", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(writer, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
