package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/internal/password"
	"github.com/mrityunjaygr8/autostrada-test/internal/request"
	"github.com/mrityunjaygr8/autostrada-test/internal/response"
	"github.com/mrityunjaygr8/autostrada-test/internal/validator"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"net/http"
	"strconv"
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
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Admin     bool                `json:"admin"`
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

	input.Validator.CheckField(input.Email != "", "email", "Email is required")
	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "email", "Must be a valid email address")
	input.Validator.CheckField(existingUser == nil, "email", "Email is already in use")

	input.Validator.CheckField(input.Password != "", "password", "Password is required")
	input.Validator.CheckField(len(input.Password) >= 8, "password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "password", "Password is too long")
	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "password", "Password is too common")

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

	err = response.JSONWithHeaders(w, http.StatusCreated, map[string]interface{}{"Data": user}, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

const DefaultPageNumber = "1"
const DefaultPageSize = "20"

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PageSize   string
		PageNumber string
		Validator  validator.Validator
	}

	var params store.UserListParams

	//pageNumber := r.URL.Query().Get("pageNumber")
	//pageSize := r.URL.Query().Get("pageSize")

	var pageSizeErr, pageNumberErr error
	input.PageSize = r.URL.Query().Get("pageSize")
	input.PageNumber = r.URL.Query().Get("pageNumber")

	if input.PageNumber == "" {
		input.PageNumber = DefaultPageNumber
	}

	if input.PageSize == "" {
		input.PageSize = DefaultPageSize
	}

	params.PageSize, pageSizeErr = strconv.Atoi(input.PageSize)
	params.PageNumber, pageNumberErr = strconv.Atoi(input.PageNumber)

	//fmt.Println(params.PageSize, params.PageSize > 0, params.PageNumber, params.PageNumber > 0)

	input.Validator.CheckField(pageSizeErr == nil, "pageSize", "pageSize must be a positive integer")
	input.Validator.CheckField(pageNumberErr == nil, "pageNumber", "pageNumber must be a positive integer")
	input.Validator.CheckField(params.PageSize > 0, "pageSize", "pageSize must be a positive integer")
	input.Validator.CheckField(params.PageNumber > 0, "pageNumber", "pageNumber must be a positive integer")

	if input.Validator.HasErrors() {
		fmt.Println(input.Validator.FieldErrors)
		app.failedValidation(w, r, input.Validator)
		return
	}
	users, err := app.store.UserList(params)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, users)
	if err != nil {
		app.serverError(w, r, err)
	}
}
