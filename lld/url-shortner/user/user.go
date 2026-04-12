package user

import (
	"time"

	"github.com/google/uuid"
)

func NewUser(name, email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
}
