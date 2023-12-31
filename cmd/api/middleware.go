package main

import (
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//defer func() {
		//	err := recover()
		//	if err != nil {
		//		app.serverError(w, r, fmt.Errorf("%s", err))
		//	}
		//}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Add("Vary", "Authorization")
		//
		//authorizationHeader := r.Header.Get("Authorization")
		//
		//if authorizationHeader != "" {
		//	headerParts := strings.Split(authorizationHeader, " ")
		//
		//	if len(headerParts) == 2 && headerParts[0] == "Bearer" {
		//		token := headerParts[1]
		//
		//		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secretKey))
		//		if err != nil {
		//			app.invalidAuthenticationToken(w, r)
		//			return
		//		}
		//
		//		if !claims.Valid(time.Now()) {
		//			app.invalidAuthenticationToken(w, r)
		//			return
		//		}
		//
		//		if claims.Issuer != app.config.baseURL {
		//			app.invalidAuthenticationToken(w, r)
		//			return
		//		}
		//
		//		if !claims.AcceptAudience(app.config.baseURL) {
		//			app.invalidAuthenticationToken(w, r)
		//			return
		//		}
		//
		//		userID, err := strconv.Atoi(claims.Subject)
		//		if err != nil {
		//			app.serverError(w, r, err)
		//			return
		//		}
		//
		//		user, err := app.store.UserRetrieve(userID)
		//		if err != nil {
		//			app.serverError(w, r, err)
		//			return
		//		}
		//
		//		if user != nil {
		//			r = contextSetAuthenticatedUser(r, user)
		//		}
		//	}
		//}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authenticatedUser := contextGetAuthenticatedUser(r)
		//
		//if authenticatedUser == nil {
		//	app.authenticationRequired(w, r)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
