package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type Expectation int

const (
	Success Expectation = iota
	FieldError
	Error
)

type args struct {
	input  string
	expect string
	code   int
}

type resp struct {
	data        map[string]any `json:"data"`
	fieldErrors map[string]any `json:"FieldErrors"`
	errors      map[string]any `json:"Error"`
}

func Test_application_createUser(t *testing.T) {
	stubStore := NewStubStore()
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "basis insert success",
			args: args{
				input: `
				{
					"Email": "msyt1969@gmail.com",
					"Password": "woowoowoo",
					"Admin": true
				}`,
				code: http.StatusCreated,
				expect: `{
					"data": {
						"email": "msyt1969@gmail.com",
						"admin": true
					}
				}`,
			},
		},
		{
			name: "basis insert bad email",
			args: args{
				input: `
				{
					"Email": "msyt1969",
					"Password": "woowoowoo",
					"Admin": true
				}`,
				code: http.StatusUnprocessableEntity,
				expect: `{
					"FieldErrors": {
						"email": "Must be a valid email address"
					}
				}`,
			},
		},
		{
			name: "basic insert no body",
			args: args{
				input:  ``,
				code:   http.StatusBadRequest,
				expect: `{"Error": "Body must not be empty"}`,
			},
		},
	}
	app := &application{
		store: stubStore,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.args.input))
			response := httptest.NewRecorder()
			app.createUser(response, request)
			require.Equal(t, tt.args.code, response.Code)
			var res map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &res)
			require.Nil(t, err)
			var expected resp
			t.Log(tt.args.expect, []byte(tt.args.expect))
			err = json.Unmarshal([]byte(tt.args.expect), &expected)
			require.Nil(t, err)
			t.Log(expected)
			//for key, val := range tt.args.expect {
			//	require.Equal(t, val, res[key])
			//}
		})
	}
}

//func Test_application_listUsers(t *testing.T) {
//	tests := []struct {
//		name  string
//		store store.GuzeiStore
//		args  args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := &application{
//				store: tt.store,
//			}
//			app.listUsers(tt.args.w, tt.args.r)
//		})
//	}
//}

func Test_application_status(t *testing.T) {
	stubStore := NewStubStore()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "basic status",
			args: args{
				code: http.StatusOK,
				expect: `{
					"Status": "OK",
				}`,
			},
		},
	}
	app := &application{
		store: stubStore,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/status", nil)
			response := httptest.NewRecorder()
			app.status(response, request)
			var res map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &res)
			require.Nil(t, err)
			//for key, val := range tt.args.expect {
			//	require.Equal(t, val, res[key])
			//}
		})
	}
}

type StubStore struct {
	userStore map[string]store.User
}

func NewStubStore() StubStore {
	userStore := make(map[string]store.User)
	return StubStore{userStore}
}

func (s StubStore) UserInsert(email, password string, id uuid.UUID, admin bool) (*store.User, error) {
	_, found := s.userStore[email]
	if found {
		return nil, store.ErrUserExists
	}
	u := store.User{
		Email:          email,
		ID:             id,
		Admin:          admin,
		Created:        time.Now(),
		HashedPassword: password,
	}
	s.userStore[email] = u
	return &u, nil
}

func (s StubStore) UserList(userListParams store.UserListParams) (*store.UsersList, error) {
	//TODO implement me
	panic("implement me")
}

func (s StubStore) UserRetrieveByEmail(email string) (*store.User, error) {
	u, found := s.userStore[email]
	if !found {
		return nil, store.ErrUserNotFound
	}
	return &u, nil
}

func (s StubStore) UserRetrieve(id uuid.UUID) (*store.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s StubStore) UserUpdatePassword(id uuid.UUID, newPassword string) error {
	//TODO implement me
	panic("implement me")
}

func (s StubStore) UserUpdateAdmin(id uuid.UUID, newAdminValue bool) error {
	//TODO implement me
	panic("implement me")
}

func (s StubStore) UserDelete(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
