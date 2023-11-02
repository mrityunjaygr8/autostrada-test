package store

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type GuzeiStore interface {
	UserInsert(email, password string, id uuid.UUID, admin bool) (*User, error)
	UserList(userListParams UserListParams) (*UsersList, error)
	UserRetrieveByEmail(email string) (*User, error)
	UserRetrieve(id uuid.UUID) (*User, error)
	UserUpdatePassword(id uuid.UUID, newPassword string) error
	UserUpdateAdmin(id uuid.UUID, newAdminValue bool) error
	UserDelete(id uuid.UUID) error
}

type UserListParams struct {
	PageNumber int
	PageSize   int
}

type User struct {
	Email          string    `json:"email"`
	ID             uuid.UUID `json:"id"`
	Admin          bool      `json:"admin"`
	Created        time.Time `json:"created"`
	HashedPassword string    `json:"-"`
}

type UsersList struct {
	Data         []User `json:"data"`
	TotalObjects int    `json:"total"`
	TotalPages   int    `json:"pages"`
}

var ErrUserExists = errors.New("user with specified email already exists")
var ErrUserNotFound = errors.New("specified user does not exists")
var ErrStoreError = errors.New("error persisting in storage")
