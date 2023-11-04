package main

import "net/http"

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
