package auth

type Username string
type Password string
type UserId string

type User struct {
	Username Username
	ID       UserId
}

type Authenticator interface {
	Login(Username, Password) (*User, error)
}
