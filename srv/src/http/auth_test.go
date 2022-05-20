package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuther(t *testing.T) {

	password := "foo"
	hashedPassword := NewPasswordHash(password)

	auther := NewAuther(map[string]string{
		"FOO": hashedPassword,
	})

	assert.False(t, auther.Allowed("BAR", password))
	assert.False(t, auther.Allowed("FOO", "bar"))
	assert.True(t, auther.Allowed("FOO", password))
}
