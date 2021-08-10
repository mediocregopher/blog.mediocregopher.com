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
	args  []string
}

// New initializes and returns a new instance of *Cfg.
func New() *Cfg {
	return &Cfg{
		FlagSet: flag.NewFlagSet("", flag.ExitOnError),
		args:    os.Args[1:],
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
	if err := c.FlagSet.Parse(c.args); err != nil {
		return err
	}

	for _, h := range c.hooks {
		if err := h(ctx); err != nil {
			return err
		}
	}

	return nil
}

// SubCmd should be called _after_ Init. Init will have consumed all arguments
// up until the first non-flag argument. This non-flag argument is a
// sub-command, and is returned by this method. This method also resets Cfg's
// internal state so that new options can be added to it.
//
// If there is no sub-command following the initial set of flags then this will
// return empty string.
func (c *Cfg) SubCmd() string {
	c.args = c.FlagSet.Args()
	if len(c.args) == 0 {
		return ""
	}

	subCmd := c.args[0]

	c.FlagSet = flag.NewFlagSet(subCmd, flag.ExitOnError)
	c.hooks = nil
	c.args = c.args[1:]

	return subCmd
}
