package chat

import (
	"encoding/hex"
	"fmt"
	"sync"

	"golang.org/x/crypto/argon2"
)

// UserID uniquely identifies an individual user who has posted a message in a
// Room.
type UserID struct {

	// Name will be the user's chosen display name.
	Name string `json:"name"`

	// Hash will be a hex string generated from a secret only the user knows.
	Hash string `json:"id"`
}

// UserIDCalculator is used to calculate UserIDs.
type UserIDCalculator struct {

	// Secret is used when calculating UserID Hash salts.
	Secret []byte

	// TimeCost, MemoryCost, and Threads are used as inputs to the Argon2id
	// algorithm which is used to generate the Hash.
	TimeCost, MemoryCost uint32
	Threads              uint8

	// HashLen specifies the number of bytes the Hash should be.
	HashLen uint32

	// Lock, if set, forces concurrent Calculate calls to occur sequentially.
	Lock *sync.Mutex
}

// NewUserIDCalculator returns a UserIDCalculator with sane defaults.
func NewUserIDCalculator(secret []byte) UserIDCalculator {
	return UserIDCalculator{
		Secret:     secret,
		TimeCost:   15,
		MemoryCost: 128 * 1024,
		Threads:    2,
		HashLen:    16,
		Lock:       new(sync.Mutex),
	}
}

// Calculate accepts a name and password and returns the calculated UserID.
func (c UserIDCalculator) Calculate(name, password string) UserID {

	input := fmt.Sprintf("%q:%q", name, password)

	hashB := argon2.IDKey(
		[]byte(input),
		c.Secret, // salt
		c.TimeCost, c.MemoryCost, c.Threads,
		c.HashLen,
	)

	return UserID{
		Name: name,
		Hash: hex.EncodeToString(hashB),
	}
}
