package main

import (
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/data"
	"github.com/jessicatarra/greenlight/internal/validator"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Authorization")

		authorizationHeader := request.Header.Get("Authorization")

		if authorizationHeader == "" {
			request = app.contextSetUser(request, data.AnonymousUser)
			next.ServeHTTP(writer, request)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(writer, request)
			default:
				app.serverErrorResponse(writer, request, err)
			}
			return
		}

		request = app.contextSetUser(request, user)

		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(writer, request)
			return
		}

		if !user.Activated {
			app.inactiveAccountResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}
		if !permissions.Include(code) {
			app.notPermittedResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	}

	return app.requireActivatedUser(fn)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Origin")

		origin := request.Header.Get("Origin")

		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {

			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {

					writer.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		next.ServeHTTP(writer, request)
	})
}
