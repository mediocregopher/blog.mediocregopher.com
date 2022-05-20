package http

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuther(t *testing.T) {

	ctx := context.Background()
	password := "foo"
	hashedPassword := NewPasswordHash(password)

	auther := NewAuther(map[string]string{
		"FOO": hashedPassword,
	}, 1*time.Millisecond)

	assert.False(t, auther.Allowed(ctx, "BAR", password))
	assert.False(t, auther.Allowed(ctx, "FOO", "bar"))
	assert.True(t, auther.Allowed(ctx, "FOO", password))
}
