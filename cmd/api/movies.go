package main

import (
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/data"
	"github.com/jessicatarra/greenlight/internal/validator"
	"net/http"
)

type createMovieRequest struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

type updateMovieRequest struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *int32   `json:"runtime"`
	Genres  []string `json:"genres"`
}

// @Summary Create a movie
// @Description Create a new movie
// @Tags Movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body createMovieRequest true "Request body"
// @Success 201 {object} data.Movie "Movie created"
// @Router /movies [post]
func (app *application) createMovieHandler(writer http.ResponseWriter, request *http.Request) {
	input := createMovieRequest{}

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

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(writer, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

// @Summary Get a movie by ID
// @Description Retrieve a movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Success 200 {object} data.Movie "Movie details"
// @Router /movies/{id} [get]
func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

// @Summary Update a movie by ID
// @Description Update an existing movie
// @Tags Movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Param request body updateMovieRequest true "Request body"
// @Success 200 {object} data.Movie "Movie updated"
// @Router /movies/{id} [put]
func (app *application) updateMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	input := updateMovieRequest{}

	err = app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

// @Summary Delete a movie by ID
// @Description Delete a movie by its ID
// @Tags Movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Movie ID"
// @Success 200
// @Router /movies/{id} [delete]
func (app *application) deleteMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}

// @Summary List movies with pagination
// @Description Fetch a list of movies with server-side pagination
// @Tags Movies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param title query string false "Movie title"
// @Param genres query []string false "Movie genres"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of movies per page"
// @Param sort query string false "Sort order"
// @Success 200 {object} []data.Movie "Movie list"
// @Router /movies [get]
func (app *application) listMoviesHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := request.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
