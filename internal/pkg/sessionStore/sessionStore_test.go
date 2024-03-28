package sessionStore

import (
	"testing"
	"time"

	"github.com/fbuedding/iota-admin/internal/pkg/auth"
)

func TestSession_IsExpired(t *testing.T) {
	type fields struct {
		Username auth.Username
		Expiry   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{
				Username: tt.fields.Username,
				Expiry:   tt.fields.Expiry,
			}
			if got := s.IsExpired(); got != tt.want {
				t.Errorf("Session.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Refresh(t *testing.T) {
	type fields struct {
		Username auth.Username
		Expiry   time.Time
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Username: tt.fields.Username,
				Expiry:   tt.fields.Expiry,
			}
			s.Refresh(tt.args.t)
		})
	}
}
