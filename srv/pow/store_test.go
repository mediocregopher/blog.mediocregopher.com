package pow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tilinna/clock"
)

func TestStore(t *testing.T) {
	clock := clock.NewMock(time.Now().Truncate(time.Hour))
	now := clock.Now()

	s := NewMemoryStore(clock)
	defer s.Close()

	seed := []byte{0}

	// mark solved should work
	err := s.MarkSolved(seed, now.Add(time.Second))
	assert.NoError(t, err)

	// mark again, should not work
	err = s.MarkSolved(seed, now.Add(time.Hour))
	assert.ErrorIs(t, err, ErrSeedSolved)

	// marking a different seed should still work
	seed2 := []byte{1}
	err = s.MarkSolved(seed2, now.Add(inMemStoreGCPeriod*2))
	assert.NoError(t, err)
	err = s.MarkSolved(seed2, now.Add(time.Hour))
	assert.ErrorIs(t, err, ErrSeedSolved)

	now = clock.Add(inMemStoreGCPeriod)
	<-s.(*inMemStore).spinLoopCh

	// first one should be markable again, second shouldnt
	err = s.MarkSolved(seed, now.Add(time.Second))
	assert.NoError(t, err)
	err = s.MarkSolved(seed2, now.Add(time.Hour))
	assert.ErrorIs(t, err, ErrSeedSolved)

	now = clock.Add(inMemStoreGCPeriod)
	<-s.(*inMemStore).spinLoopCh

	// now both should be expired
	err = s.MarkSolved(seed, now.Add(time.Second))
	assert.NoError(t, err)
	err = s.MarkSolved(seed2, now.Add(time.Second))
	assert.NoError(t, err)
}
