package pow

import (
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tilinna/clock"
)

func TestChallengeParams(t *testing.T) {
	tests := []challengeParams{
		{},
		{
			Target:    1,
			ExpiresAt: 3,
		},
		{
			Target:    2,
			ExpiresAt: -10,
			Random:    []byte{0, 1, 2},
		},
		{
			Random: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	t.Run("marshal_unmarshal", func(t *testing.T) {
		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				b, err := test.MarshalBinary()
				assert.NoError(t, err)

				var c2 challengeParams
				assert.NoError(t, c2.UnmarshalBinary(b))
				assert.Equal(t, test, c2)

				b2, err := c2.MarshalBinary()
				assert.NoError(t, err)
				assert.Equal(t, b, b2)
			})
		}
	})

	secret := []byte("shhh")

	t.Run("to_from_seed", func(t *testing.T) {

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				seed, err := newSeed(test, secret)
				assert.NoError(t, err)

				// generating seed should be deterministic
				seed2, err := newSeed(test, secret)
				assert.NoError(t, err)
				assert.Equal(t, seed, seed2)

				c, err := challengeParamsFromSeed(seed, secret)
				assert.NoError(t, err)
				assert.Equal(t, test, c)
			})
		}
	})

	t.Run("malformed_seed", func(t *testing.T) {
		tests := []string{
			"",
			"01",
			"0000",
			"00374a1ad84d6b7a93e68042c1f850cbb100000000000000000000000000000102030405060708A0", // changed one byte from a good seed
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				seed, err := hex.DecodeString(test)
				if err != nil {
					panic(err)
				}

				_, err = challengeParamsFromSeed(seed, secret)
				assert.ErrorIs(t, errMalformedSeed, err)
			})
		}
	})
}

func TestManager(t *testing.T) {
	clock := clock.NewMock(time.Now().Truncate(time.Hour))

	store := NewMemoryStore(clock)
	defer store.Close()

	mgr := NewManager(ManagerParams{
		Clock:            clock,
		Store:            store,
		Secret:           []byte("shhhh"),
		Target:           0x00FFFFFF,
		ChallengeTimeout: 1 * time.Second,
	})

	{
		c := mgr.NewChallenge()
		solution := Solve(c)
		assert.NoError(t, mgr.CheckSolution(c.Seed, solution))

		// doing again should fail, the seed should already be marked as solved
		assert.ErrorIs(t, mgr.CheckSolution(c.Seed, solution), ErrSeedSolved)
	}

	{
		c := mgr.NewChallenge()
		solution := Solve(c)
		clock.Add(2 * time.Second)
		assert.ErrorIs(t, mgr.CheckSolution(c.Seed, solution), ErrExpiredSolution)
	}

}
