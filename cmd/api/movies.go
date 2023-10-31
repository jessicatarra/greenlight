package main

import (
	"fmt"
	"github.com/jessicatarra/greenlight/internal/data"
	"net/http"
	"time"
)

func (app *application) createMovieHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "create movie handler")
}

func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		http.NotFound(writer, request)
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

	err = app.writeJSON(writer, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(writer, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
