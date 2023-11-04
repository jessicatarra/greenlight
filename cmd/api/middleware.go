package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.logger.PrintInfo("request", map[string]string{
			"request_remote_addr":     request.RemoteAddr,
			"request_proto":           request.Proto,
			"request_method":          request.Method,
			"request_url_request_uri": request.URL.RequestURI(),
		})

		next.ServeHTTP(writer, request)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				writer.Header().Set("Connection", "close")

				app.serverErrorResponse(writer, request, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
