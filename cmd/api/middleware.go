package main

import (
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	limiter := rate.NewLimiter(10, 40)

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
			request = app.contextSetUser(request, database.AnonymousUser)
			next.ServeHTTP(writer, request)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		if claims.Issuer != "greenlight.tarralva.com" {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		if !claims.AcceptAudience("greenlight.tarralva.com") {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}
		user, err := app.models.Users.Get(userID)

		if err != nil {
			switch {
			case errors.Is(err, database.ErrRecordNotFound):
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

		writer.Header().Add("Vary", "Access-Control-Request-Method")

		origin := request.Header.Get("Origin")

		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					writer.Header().Set("Access-Control-Allow-Origin", origin)

					if request.Method == http.MethodOptions && request.Header.Get("Access-Control-Request-Method") != "" {

						writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						writer.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		next.ServeHTTP(writer, request)
	})
}
