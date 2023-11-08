package auth

type DebugAuth struct {
}

func NewDebugAuth() *DebugAuth {
	return &DebugAuth{}
}

func (d DebugAuth) Login(u Username, p Password) (*User, error) {
	user := &User{u, UserId("uuid:1234")}
	return user, nil
}
