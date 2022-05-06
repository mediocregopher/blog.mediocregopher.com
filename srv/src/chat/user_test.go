package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIDCalculator(t *testing.T) {

	const name, password = "name", "password"

	c := NewUserIDCalculator([]byte("foo"))

	// calculating with same params twice should result in same UserID
	userID := c.Calculate(name, password)
	assert.Equal(t, userID, c.Calculate(name, password))

	// changing either name or password should result in a different Hash
	assert.NotEqual(t, userID.Hash, c.Calculate(name+"!", password).Hash)
	assert.NotEqual(t, userID.Hash, c.Calculate(name, password+"!").Hash)

	// changing the secret should change the UserID
	c.Secret = []byte("bar")
	assert.NotEqual(t, userID, c.Calculate(name, password))
}
