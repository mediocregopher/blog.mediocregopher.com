package cfg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringVar(t *testing.T) {

	cfg := New(Params{
		Args: []string{"--foo=CLI"},
		Env:  map[string]string{"FOO": "ENV", "BAR": "ENV"},
	})

	var foo, bar, baz string

	cfg.StringVar(&foo, "foo", "DEF", "")
	cfg.StringVar(&bar, "bar", "DEF", "")
	cfg.StringVar(&baz, "baz", "DEF", "")

	assert.NoError(t, cfg.Init(context.Background()))
	assert.Equal(t, "CLI", foo)
	assert.Equal(t, "ENV", bar)
	assert.Equal(t, "DEF", baz)
}

func TestIntVar(t *testing.T) {

	cfg := New(Params{
		Args: []string{"--foo=111"},
		Env:  map[string]string{"FOO": "222", "BAR": "222"},
	})

	var foo, bar, baz int

	cfg.IntVar(&foo, "foo", 333, "")
	cfg.IntVar(&bar, "bar", 333, "")
	cfg.IntVar(&baz, "baz", 333, "")

	assert.NoError(t, cfg.Init(context.Background()))
	assert.Equal(t, 111, foo)
	assert.Equal(t, 222, bar)
	assert.Equal(t, 333, baz)
}
