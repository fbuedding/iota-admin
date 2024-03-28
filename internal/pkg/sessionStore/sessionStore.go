package sessionStore

import (
	"time"

	"github.com/fbuedding/iota-admin/internal/pkg/auth"
)

type Session struct {
	Username auth.Username
	Expiry   time.Time
}
type SessionToken string

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func (s *Session) Refresh(t time.Time) {
	s.Expiry = t
}

type SessionStore interface {
	Add(*Session) (SessionToken, error)
	Get(SessionToken) (*Session, error)
	Exists(SessionToken) bool
	Remove(SessionToken)
}
