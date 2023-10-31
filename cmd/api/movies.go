package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(writer, "show movie handler %d", id)
}
