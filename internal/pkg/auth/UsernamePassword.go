package auth

import (
	"github.com/fbuedding/iota-admin/internal/globals"
)

type UsernamePassword struct{}

func (up UsernamePassword) Authenticate(username Username, password Password) (*User, error) {
	if username == Username(globals.Conf.Username) && password == Password(globals.Conf.Password) {
		return &User{
			Username: username,
			ID:       "environment_vars",
			Role:     "Admin",
		}, nil
	}
	return nil, ErrAuthFailed
}

func NewUsernamePasswordAuth() *UsernamePassword {
	return &UsernamePassword{}
}
