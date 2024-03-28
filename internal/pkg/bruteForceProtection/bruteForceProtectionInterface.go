package bruteforceprotection

import "github.com/fbuedding/iota-admin/internal/pkg/auth"

type BrutForceProtection interface {
	IsBlocked(auth.Username) bool
	Hit(auth.Username)
	Delete(auth.Username)
}
