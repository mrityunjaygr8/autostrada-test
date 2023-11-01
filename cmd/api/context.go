package main

import (
	"context"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*store.User)
	if !ok {
		return nil
	}

	return user
}
