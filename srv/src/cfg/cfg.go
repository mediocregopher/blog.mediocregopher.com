// Package cfg implements a simple wrapper around go's flag package, in order to
// implement initialization hooks.
package cfg

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Cfger is a component which can be used with Cfg to setup its initialization.
type Cfger interface {
	SetupCfg(*Cfg)
}

// Params are used to initialize a Cfg instance.
type Params struct {

	// Args are the command line arguments, excluding the command-name.
	//
	// Defaults to os.Args[1:]
	Args []string

	// Env is the process's environment variables.
	//
	// Defaults to the real environment variables.
	Env map[string]string

	// EnvPrefix indicates a string to prefix to all environment variable names
	// that Cfg will read. Will be automatically suffixed with a "_" if given.
	EnvPrefix string
}

func (p Params) withDefaults() Params {

	if p.Args == nil {
		p.Args = os.Args[1:]
	}

	if p.Env == nil {

		p.Env = map[string]string{}

		for _, envVar := range os.Environ() {

			parts := strings.SplitN(envVar, "=", 2)

			if len(parts) < 2 {
				panic(fmt.Sprintf("envVar %q returned from os.Environ() somehow", envVar))
			}

			p.Env[parts[0]] = parts[1]
		}
	}

	if p.EnvPrefix != "" {
		p.EnvPrefix = strings.TrimSuffix(p.EnvPrefix, "_") + "_"
	}

	return p
}

// Cfg is a wrapper around the stdlib's FlagSet and a set of initialization
// hooks.
type Cfg struct {
	params  Params
	flagSet *flag.FlagSet

	hooks []func(ctx context.Context) error
}

// New initializes and returns a new instance of *Cfg.
func New(params Params) *Cfg {

	params = params.withDefaults()

	return &Cfg{
		params:  params,
		flagSet: flag.NewFlagSet("", flag.ExitOnError),
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
	if err := c.flagSet.Parse(c.params.Args); err != nil {
		return err
	}

	for _, h := range c.hooks {
		if err := h(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cfg) envifyName(name string) string {
	name = c.params.EnvPrefix + name
	name = strings.Replace(name, "-", "_", -1)
	name = strings.ToUpper(name)
	return name
}

func envifyUsage(envName, usage string) string {
	return fmt.Sprintf("%s (overrides %s)", usage, envName)
}

// StringVar is equivalent to flag.FlagSet's StringVar method, but will
// additionally set up an environment variable for the parameter.
func (c *Cfg) StringVar(p *string, name, value, usage string) {

	envName := c.envifyName(name)

	c.flagSet.StringVar(p, name, value, envifyUsage(envName, usage))

	if val := c.params.Env[envName]; val != "" {
		*p = val
	}
}

// Args returns a pointer which will be filled with the process's positional
// arguments after Init is called. The positional arguments are all CLI
// arguments starting with the first non-flag argument.
//
// The usage argument should describe what these arguments are, and its notation
// should indicate if they are optional or variadic. For example:
//
//	// optional variadic
//	"[names...]"
//
//	// required single args
//	"<something> <something else>"
//
//	// Mixed
//	"<foo> <bar> [baz] [other...]"
//
func (c *Cfg) Args(usage string) *[]string {

	args := new([]string)

	c.flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE [flags...] %s\n", usage)
		fmt.Fprintf(os.Stderr, "\nFLAGS\n\n")
		c.flagSet.PrintDefaults()
	}

	c.OnInit(func(ctx context.Context) error {
		*args = c.flagSet.Args()
		return nil
	})

	return args
}

// String is equivalent to flag.FlagSet's String method, but will additionally
// set up an environment variable for the parameter.
func (c *Cfg) String(name, value, usage string) *string {
	p := new(string)
	c.StringVar(p, name, value, usage)
	return p
}

// IntVar is equivalent to flag.FlagSet's IntVar method, but will additionally
// set up an environment variable for the parameter.
func (c *Cfg) IntVar(p *int, name string, value int, usage string) {

	envName := c.envifyName(name)

	c.flagSet.IntVar(p, name, value, envifyUsage(envName, usage))

	// if we can't parse the envvar now then just hold onto the error until
	// Init, otherwise we'd have to panic here and that'd be ugly.
	var err error

	if valStr := c.params.Env[envName]; valStr != "" {

		var val int
		val, err = strconv.Atoi(valStr)

		if err != nil {
			err = fmt.Errorf(
				"parsing envvar %q with value %q: %w",
				envName, valStr, err,
			)

		} else {
			*p = val
		}
	}

	c.OnInit(func(context.Context) error { return err })
}

// Int is equivalent to flag.FlagSet's Int method, but will additionally set up
// an environment variable for the parameter.
func (c *Cfg) Int(name string, value int, usage string) *int {
	p := new(int)
	c.IntVar(p, name, value, usage)
	return p
}

// BoolVar is equivalent to flag.FlagSet's BoolVar method, but will additionally
// set up an environment variable for the parameter.
func (c *Cfg) BoolVar(p *bool, name string, value bool, usage string) {

	envName := c.envifyName(name)

	c.flagSet.BoolVar(p, name, value, envifyUsage(envName, usage))

	if valStr := c.params.Env[envName]; valStr != "" {
		*p = valStr != "" && valStr != "0" && valStr != "false"
	}
}

// Bool is equivalent to flag.FlagSet's Bool method, but will additionally set
// up an environment variable for the parameter.
func (c *Cfg) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	c.BoolVar(p, name, value, usage)
	return p
}

// SubCmd should be called _after_ Init. Init will have consumed all arguments
// up until the first non-flag argument. This non-flag argument is a
// sub-command, and is returned by this method. This method also resets Cfg's
// internal state so that new options can be added to it.
//
// If there is no sub-command following the initial set of flags then this will
// return empty string.
func (c *Cfg) SubCmd() string {
	c.params.Args = c.flagSet.Args()
	if len(c.params.Args) == 0 {
		return ""
	}

	subCmd := c.params.Args[0]

	c.flagSet = flag.NewFlagSet(subCmd, flag.ExitOnError)
	c.hooks = nil
	c.params.Args = c.params.Args[1:]

	return subCmd
}
