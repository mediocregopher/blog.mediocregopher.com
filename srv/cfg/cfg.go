// Package cfg implements a simple wrapper around go's flag package, in order to
// implement initialization hooks.
package cfg

import (
	"context"
	"flag"
	"os"
)

// Cfger is a component which can be used with Cfg to setup its initialization.
type Cfger interface {
	SetupCfg(*Cfg)
}

// Cfg is a wrapper around the stdlib's FlagSet and a set of initialization
// hooks.
type Cfg struct {
	*flag.FlagSet

	hooks []func(ctx context.Context) error
}

// New initializes and returns a new instance of *Cfg.
func New() *Cfg {
	return &Cfg{
		FlagSet: flag.NewFlagSet("", flag.ExitOnError),
	}
}

// OnInit appends the given callback to the sequence of hooks which will run on
// a call to Init.
func (c *Cfg) OnInit(cb func(context.Context) error) {
	c.hooks = append(c.hooks, cb)
}

// Init runs all hooks registered using OnInit, in the same order OnInit was
// called. If one returns an error that error is returned and no further hooks
// are run.
func (c *Cfg) Init(ctx context.Context) error {
	if err := c.FlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	for _, h := range c.hooks {
		if err := h(ctx); err != nil {
			return err
		}
	}

	return nil
}
