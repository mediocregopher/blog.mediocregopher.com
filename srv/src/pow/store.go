package pow

import (
	"errors"
	"sync"
	"time"

	"github.com/tilinna/clock"
)

// ErrSeedSolved is used to indicate a seed has already been solved.
var ErrSeedSolved = errors.New("seed already solved")

// Store is used to track information related to proof-of-work challenges and
// solutions.
type Store interface {

	// MarkSolved will return ErrSeedSolved if the seed was already marked. The
	// seed will be cleared from the Store once expiresAt is reached.
	MarkSolved(seed []byte, expiresAt time.Time) error

	Close() error
}

type inMemStore struct {
	clock clock.Clock

	m          map[string]time.Time
	l          sync.Mutex
	closeCh    chan struct{}
	spinLoopCh chan struct{} // only used by tests
}

const inMemStoreGCPeriod = 5 * time.Second

// NewMemoryStore initializes and returns an in-memory Store implementation.
func NewMemoryStore(clock clock.Clock) Store {
	s := &inMemStore{
		clock:      clock,
		m:          map[string]time.Time{},
		closeCh:    make(chan struct{}),
		spinLoopCh: make(chan struct{}, 1),
	}
	go s.spin(s.clock.NewTicker(inMemStoreGCPeriod))
	return s
}

func (s *inMemStore) spin(ticker *clock.Ticker) {
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := s.clock.Now()

			s.l.Lock()
			for seed, expiresAt := range s.m {
				if !now.Before(expiresAt) {
					delete(s.m, seed)
				}
			}
			s.l.Unlock()

		case <-s.closeCh:
			return
		}

		select {
		case s.spinLoopCh <- struct{}{}:
		default:
		}
	}
}

func (s *inMemStore) MarkSolved(seed []byte, expiresAt time.Time) error {
	seedStr := string(seed)

	s.l.Lock()
	defer s.l.Unlock()

	if _, ok := s.m[seedStr]; ok {
		return ErrSeedSolved
	}

	s.m[seedStr] = expiresAt
	return nil
}

func (s *inMemStore) Close() error {
	close(s.closeCh)
	return nil
}
