// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID             uuid.UUID
	Created        pgtype.Timestamptz
	Email          string
	Admin          bool
	HashedPassword string
}
