package main

import (
	"context"
	"github.com/jessicatarra/greenlight/internal/database"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *database.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *database.User {
	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
