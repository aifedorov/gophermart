package user

import (
	"github.com/google/uuid"
	"time"
)

type ID string
type User struct {
	ID           ID
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(login string, passwordHash string) *User {
	now := time.Now()
	return &User{
		ID:           ID(uuid.NewString()),
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
