package auth

import "fmt"

type (
	Username string
	Password string
	UserId   string
)

type User struct {
	Username Username
	ID       UserId
	Role     string
}

var ErrAuthFailed = fmt.Errorf("Authentcation failed!")

type Authenticator interface {
	Authenticate(Username, Password) (*User, error)
}
