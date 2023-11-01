package main

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/internal/password"
	"github.com/mrityunjaygr8/autostrada-test/internal/request"
	"github.com/mrityunjaygr8/autostrada-test/internal/response"
	"github.com/mrityunjaygr8/autostrada-test/internal/validator"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"net/http"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Admin     bool                `json:"Admin"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	existingUser, err := app.store.UserRetrieveByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			break
		default:
			app.serverError(w, r, err)
			return
		}
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
	input.Validator.CheckField(existingUser == nil, "Email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	id := uuid.New()
	user, err := app.store.UserInsert(input.Email, hashedPassword, id, input.Admin)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSONWithHeaders(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string              `json:"Email"`
		Password  string              `json:"Password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.store.UserRetrieveByEmail(input.Email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(user != nil, "Email", "Email address could not be found")

	if user != nil {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	//var claims jwt.Claims
	//claims.Subject = strconv.Itoa(user.ID)
	//
	//expiry := time.Now().Add(24 * time.Hour)
	//claims.Issued = jwt.NewNumericTime(time.Now())
	//claims.NotBefore = jwt.NewNumericTime(time.Now())
	//claims.Expires = jwt.NewNumericTime(expiry)
	//
	//claims.Issuer = app.config.baseURL
	//claims.Audiences = []string{app.config.baseURL}
	//
	//jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	//if err != nil {
	//	app.serverError(w, r, err)
	//	return
	//}
	//
	//data := map[string]string{
	//	"AuthenticationToken":       string(jwtBytes),
	//	"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	//}
	//
	//err = response.JSON(w, http.StatusOK, data)
	//if err != nil {
	//	app.serverError(w, r, err)
	//}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
