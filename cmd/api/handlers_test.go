package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type resp struct {
	Data        map[string]any `json:"Data,omitempty"`
	FieldErrors map[string]any `json:"FieldErrors,omitempty"`
	Errors      string         `json:"Error,omitempty"`
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

type createUserProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

func TestCreateUser(t *testing.T) {
	stubStore := NewStubStore()
	app := &application{
		store: &stubStore,
	}
	t.Run("TestCreateUser happy path", func(t *testing.T) {
		data := createUserProps{
			Email:    "msyt@gmail.com",
			Password: "qweqweqwe",
			Admin:    false,
		}
		bData, err := json.Marshal(data)
		require.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
		response := httptest.NewRecorder()

		app.createUser(response, request)

		require.Equal(t, http.StatusCreated, response.Code)
		var res resp
		err = json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, data.Admin, res.Data["admin"])
		require.Equal(t, data.Email, res.Data["email"])
		require.True(t, IsValidUUID(res.Data["id"].(string)))
	})

	t.Run("TestCreateUser error case", func(t *testing.T) {
		data := createUserProps{
			Email:    "msyt",
			Password: "qweqweqwe",
			Admin:    false,
		}
		bData, err := json.Marshal(data)
		require.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
		response := httptest.NewRecorder()

		app.createUser(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		var res resp
		err = json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, "Must be a valid email address", res.FieldErrors["email"])
	})
	t.Run("TestCreateUser password too short", func(t *testing.T) {
		data := createUserProps{
			Email:    "msyt@gmail.com",
			Password: "qwe",
			Admin:    false,
		}
		bData, err := json.Marshal(data)
		require.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
		response := httptest.NewRecorder()

		app.createUser(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		var res resp
		err = json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, "Password is too short", res.FieldErrors["password"])
	})
	t.Run("TestCreateUser user already exists", func(t *testing.T) {
		dataOld := createUserProps{
			Email:    "msyt@gmail.com",
			Password: "qweqweqwe",
			Admin:    false,
		}
		bDataOld, err := json.Marshal(dataOld)
		require.Nil(t, err)
		requestOld := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bDataOld))
		responseOld := httptest.NewRecorder()

		app.createUser(responseOld, requestOld)
		data := createUserProps{
			Email:    "msyt@gmail.com",
			Password: "qweqweqwe",
			Admin:    false,
		}
		bData, err := json.Marshal(data)
		require.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
		response := httptest.NewRecorder()

		app.createUser(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		var res resp
		err = json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, "Email is already in use", res.FieldErrors["email"])
	})

	t.Run("CreateUser bad JSON", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"asd"`))
		response := httptest.NewRecorder()

		app.createUser(response, request)

		out := `{
	"Error": "Body contains badly-formed JSON"
}
`
		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, out, response.Body.String())

	})
}

func TestListUsers(t *testing.T) {
	t.Run("ListUsers basic case", func(t *testing.T) {
		stubStore := NewStubStore()
		app := &application{
			store: &stubStore,
		}
		for x := 0; x < 34; x++ {
			data := createUserProps{
				Email:    fmt.Sprintf("msyt_%d@gmail.com", x),
				Password: "qweqweqwe",
				Admin:    false,
			}
			bData, err := json.Marshal(data)
			require.Nil(t, err)
			requestOld := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
			responseOld := httptest.NewRecorder()

			app.createUser(responseOld, requestOld)

			require.Equal(t, http.StatusCreated, responseOld.Code)

		}
		request := httptest.NewRequest(http.MethodGet, "/users", nil)
		response := httptest.NewRecorder()

		app.listUsers(response, request)

		var res store.UsersList
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, 20, res.PageSize)
		require.Equal(t, 1, res.Page)
		require.Equal(t, 2, res.TotalPages)
		require.Equal(t, 34, res.TotalObjects)
	})
	t.Run("ListUsers query params", func(t *testing.T) {
		stubStore := NewStubStore()
		app := &application{
			store: &stubStore,
		}
		for x := 0; x < 4; x++ {
			data := createUserProps{
				Email:    fmt.Sprintf("msyt_%d@gmail.com", x),
				Password: "qweqweqwe",
				Admin:    false,
			}
			bData, err := json.Marshal(data)
			require.Nil(t, err)
			requestOld := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bData))
			responseOld := httptest.NewRecorder()

			app.createUser(responseOld, requestOld)

			require.Equal(t, http.StatusCreated, responseOld.Code)

		}
		pageSize := 2
		pageNumber := 2
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users?pageNumber=%d&pageSize=%d", pageNumber, pageSize), nil)
		response := httptest.NewRecorder()

		app.listUsers(response, request)

		var res store.UsersList
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.Nil(t, err)
		require.Equal(t, pageSize, res.PageSize)
		require.Equal(t, pageNumber, res.Page)
		require.Equal(t, 4, res.TotalObjects)
		require.Equal(t, int(math.Ceil(float64(4)/float64(pageSize))), res.TotalPages)
		require.Equal(t, 2, len(res.Data))
	})
	t.Run("ListUsers query params bad page size", func(t *testing.T) {
		stubStore := NewStubStore()
		app := &application{
			store: &stubStore,
		}
		request := httptest.NewRequest(http.MethodGet, "/users?pageSize=qwe", nil)
		response := httptest.NewRecorder()

		app.listUsers(response, request)
		out := `{
	"FieldErrors": {
		"pageSize": "pageSize must be a positive integer"
	}
}
`

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Equal(t, out, response.Body.String())
	})
	t.Run("ListUsers query params bad page number", func(t *testing.T) {
		stubStore := NewStubStore()
		app := &application{
			store: &stubStore,
		}
		request := httptest.NewRequest(http.MethodGet, "/users?pageNumber=-123", nil)
		response := httptest.NewRecorder()

		app.listUsers(response, request)
		out := `{
	"FieldErrors": {
		"pageNumber": "pageNumber must be a positive integer"
	}
}
`

		//t.Log(response.Body.String())
		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Equal(t, out, response.Body.String())
	})
}

func TestStatus(t *testing.T) {
	stubStore := NewStubStore()
	app := &application{
		store: &stubStore,
	}

	t.Run("Status check", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/status", nil)
		response := httptest.NewRecorder()

		app.status(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		out := `{
	"Status": "OK"
}
`
		require.Equal(t, out, response.Body.String())

	})
}

type StubStore struct {
	userStore []store.User
}

func NewStubStore() StubStore {
	userStore := make([]store.User, 0)
	return StubStore{userStore}
}

func (s *StubStore) UserInsert(email, password string, id uuid.UUID, admin bool) (*store.User, error) {
	u := store.User{
		Email:          email,
		ID:             id,
		Admin:          admin,
		Created:        time.Now(),
		HashedPassword: password,
	}
	s.userStore = append(s.userStore, u)
	return &u, nil
}

func (s *StubStore) UserList(userListParams store.UserListParams) (*store.UsersList, error) {
	users := make([]store.User, 0)
	count := len(s.userStore)
	start := userListParams.PageSize * (userListParams.PageNumber - 1)
	for _, val := range s.userStore[start : start+userListParams.PageSize] {
		users = append(users, val)
	}
	res := store.UsersList{
		Data:         users,
		TotalObjects: count,
		TotalPages:   int(math.Ceil(float64(count) / float64(userListParams.PageSize))),
		Page:         userListParams.PageNumber,
		PageSize:     userListParams.PageSize,
	}

	return &res, nil
}

func (s *StubStore) UserRetrieveByEmail(email string) (*store.User, error) {
	for _, item := range s.userStore {
		if item.Email == email {
			return &item, nil
		}
	}
	return nil, store.ErrUserNotFound
}

func (s *StubStore) UserRetrieve(id uuid.UUID) (*store.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StubStore) UserUpdatePassword(id uuid.UUID, newPassword string) error {
	//TODO implement me
	panic("implement me")
}

func (s *StubStore) UserUpdateAdmin(id uuid.UUID, newAdminValue bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *StubStore) UserDelete(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
