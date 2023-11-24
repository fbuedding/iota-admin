package auth

import "github.com/rs/zerolog/log"

type DebugAuth struct {
}

func NewDebugAuth() *DebugAuth {
	return &DebugAuth{}
}

func (d DebugAuth) Login(u Username, p Password) (*User, error) {
	if u == "" {
		log.Debug().Msg("Username not set")
	}
	if p == "" {
		log.Debug().Msg("Password not set")
	}
	user := &User{u, UserId("uuid:1234")}
	return user, nil
}
