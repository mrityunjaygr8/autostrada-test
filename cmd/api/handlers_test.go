package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/internal/smtp"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

//func Test_application_createAuthenticationToken(t *testing.T) {
//	type fields struct {
//		config config
//		store  store.GuzeiStore
//		logger *slog.Logger
//		mailer *smtp.Mailer
//		wg     sync.WaitGroup
//	}
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			app := &application{
//				config: tt.fields.config,
//				store:  tt.fields.store,
//				logger: tt.fields.logger,
//				mailer: tt.fields.mailer,
//				wg:     tt.fields.wg,
//			}
//			app.createAuthenticationToken(tt.args.w, tt.args.r)
//		})
//	}
//}

func Test_application_createUser(t *testing.T) {
	type Expection int
	const (
		Success Expection = iota
		FieldError
		Error
	)
	type fields struct {
		config config
		store  store.GuzeiStore
		logger *slog.Logger
		mailer *smtp.Mailer
		wg     sync.WaitGroup
	}
	type args struct {
		w           *httptest.ResponseRecorder
		r           *http.Request
		expectation Expection
	}

	type errorMessage struct {
		Error string
	}
	store := NewStubStore()
	f := fields{
		store: store,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{

			name:   "basis insert success",
			fields: f,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`
				{
					"Email": "msyt1969@gmail.com",
					"Password": "woowoowoo",
					"Admin": true
				}`)),
				expectation: Error,
			},
		},
		{
			name:   "basis insert no body",
			fields: f,
			args: args{
				w:           httptest.NewRecorder(),
				r:           httptest.NewRequest(http.MethodPost, "/users", nil),
				expectation: Error,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				config: tt.fields.config,
				store:  tt.fields.store,
				logger: tt.fields.logger,
				mailer: tt.fields.mailer,
				wg:     tt.fields.wg,
			}
			app.createUser(tt.args.w, tt.args.r)
			if tt.args.expectation == Error {
				var e errorMessage
				err := json.NewDecoder(tt.args.w.Body).Decode(&e)
				require.Nil(t, err)
				require.Equal(t, "Body must not be empty", e.Error)
				require.Equal(t, http.StatusBadRequest, tt.args.w.Code)

			}
		})
	}
}

func Test_application_listUsers(t *testing.T) {
	type fields struct {
		config config
		store  store.GuzeiStore
		logger *slog.Logger
		mailer *smtp.Mailer
		wg     sync.WaitGroup
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				config: tt.fields.config,
				store:  tt.fields.store,
				logger: tt.fields.logger,
				mailer: tt.fields.mailer,
				wg:     tt.fields.wg,
			}
			app.listUsers(tt.args.w, tt.args.r)
		})
	}
}

func Test_application_status(t *testing.T) {
	type fields struct {
		config config
		store  store.GuzeiStore
		logger *slog.Logger
		mailer *smtp.Mailer
		wg     sync.WaitGroup
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	store := NewStubStore()

	f := fields{
		store: store,
	}
	type res struct {
		Status string `json:"Status"`
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "basic status",
			fields: f,
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/status", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				config: tt.fields.config,
				store:  tt.fields.store,
				logger: tt.fields.logger,
				mailer: tt.fields.mailer,
				wg:     tt.fields.wg,
			}
			app.status(tt.args.w, tt.args.r)
			require.NotNil(t, tt.args.w)
			require.Equal(t, http.StatusOK, tt.args.w.Code)
			require.Equal(t, "application/json", tt.args.w.Header().Get("Content-Type"))
			var m res
			err := json.NewDecoder(tt.args.w.Body).Decode(&m)
			require.Nil(t, err)
			require.Equal(t, "OK", m.Status)
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
