package sessionStore

import (
	"errors"

	"github.com/google/uuid"
)

var Sessions = map[SessionToken]*Session{}

type InMemory struct{}

func (im InMemory) Add(s *Session) (SessionToken, error) {
	sessionToken := SessionToken(uuid.NewString())
	Sessions[sessionToken] = s
	return sessionToken, nil
}

func (im InMemory) Get(st SessionToken) (*Session, error) {
	session, exists := Sessions[st]
	if !exists {
		return nil, errors.New("Session not found!")
	}

	return session, nil
}

func (im InMemory) Exists(st SessionToken) bool {
	_, exists := Sessions[st]
	return exists
}

func (im InMemory) Remove(st SessionToken) {
	delete(Sessions, st)
}

func NewInMemory() *InMemory {
	return &InMemory{}
}
