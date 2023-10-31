package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type envelope map[string]interface{}

func (app *application) readIDParam(request *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(request.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJSON(writer http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	for key, value := range headers {
		writer.Header()[key] = value
	}

	json = append(json, '\n')

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(json)

	return nil
}
