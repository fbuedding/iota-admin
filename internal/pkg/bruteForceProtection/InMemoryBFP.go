package bruteforceprotection

import (
	"math"
	"sync"
	"time"

	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	"github.com/rs/zerolog/log"
)

type (
	inMemoryBruteForceProtection struct {
		maxTries int
		entries  sync.Map
	}
	inMemoryEntry struct {
		Count        int
		TimesBlocked int
		BlockedTime  time.Time
	}
)

var defaultEntry = inMemoryEntry{
	Count:        1,
	TimesBlocked: 0,
	BlockedTime:  time.Time{},
}

func (im *inMemoryBruteForceProtection) IsBlocked(usrn auth.Username) bool {
	result, ok := im.entries.Load(usrn)
	if !ok {
		return false
	}
	switch v := result.(type) {
	case inMemoryEntry:
		return v.BlockedTime.After(time.Now())
	default:
		log.Panic().Msg("Unexpected entry in InMemoryBruteForceProtection")
	}
	return false
}

func (im *inMemoryBruteForceProtection) Hit(usrn auth.Username) {
	// Should not happen, but if it does do nothing since user is blocked
	if im.IsBlocked(usrn) {
		return
	}
	entry, loaded := im.entries.LoadOrStore(usrn, defaultEntry)
	if loaded {
		switch v := entry.(type) {
		case inMemoryEntry:
			v.Count++
			if v.Count >= im.maxTries {
				v.BlockedTime = time.Now().Add(time.Duration(math.Pow(3, float64(v.TimesBlocked))) * 30 * time.Second)
				v.TimesBlocked++
				log.Info().Str("user", string(usrn)).Any("entry", v).Msg("User has been blocked")
			}
			im.entries.Swap(usrn, v)
		default:
			log.Panic().Msg("Unexpected entry in InMemoryBruteForceProtection")
		}
	}
}

func (im *inMemoryBruteForceProtection) Delete(usrn auth.Username) {
	im.entries.Delete(usrn)
}

func NewInMemory(tries int) *inMemoryBruteForceProtection {
	return &inMemoryBruteForceProtection{
		maxTries: tries,
		entries:  sync.Map{},
	}
}
