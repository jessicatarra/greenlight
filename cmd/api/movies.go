package main

import (
	"fmt"
	"github.com/jessicatarra/greenlight/internal/data"
	"github.com/jessicatarra/greenlight/internal/validator"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	data.ValidateMovie(v, movie)

	if !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	fmt.Fprintf(writer, "%+v\n", input)
}

func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(writer, http.StatusOK,
		envelope{
			"movie": movie,
		}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
