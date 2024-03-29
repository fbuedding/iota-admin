package sessionStore

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var Sessions sync.Map
var ErrUnkownType = errors.New("Unknown type")

type InMemory struct{}

func (im InMemory) Add(s *Session) (SessionToken, error) {
	sessionToken := SessionToken(uuid.NewString())
	Sessions.Store(sessionToken, s)
	return sessionToken, nil
}

func (im InMemory) Get(st SessionToken) (*Session, error) {
	session, exists := Sessions.Load(st)
	if !exists {
		return nil, errors.New("Session not found!")
	}

	switch v := session.(type) {

	case *Session:
		return v, nil
	default:
		return nil, ErrUnkownType
	}
}

func (im InMemory) Exists(st SessionToken) bool {
	_, exists := Sessions.Load(st)
	return exists
}

func (im InMemory) Remove(st SessionToken) {
	Sessions.Delete(st)
}

func NewInMemory() *InMemory {
	return &InMemory{}
}
