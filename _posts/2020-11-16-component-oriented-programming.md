---
title: >-
    Component Oriented Programming
description: >-
    A concise description of.
---

[A previous post in this
blog](2019-08-02-program-structure-and-composability.html) focused on a
framework developed to make designing component-based programs easier. In
retrospect pattern/framework proposed was over-engineered; this post attempts to
present the same ideas but in a more distilled form, as a simple programming
pattern and without the unnecessary framework.

Nothing in this post will be revelatory; it's surely all been said before. But
hopefully the form it takes here will be useful to someone, as it would have
been useful to myself when I first learned to program.

## Axioms

For the sake of brevity let's assume the following: within the context of
single-process (_not_ the same as single-threaded), non-graphical programs the
following may be said:

1. A program may be thought of as a black-box with certain input and output
   methods. It is the programmer's task to construct a program such that
   specific inputs yield specific desired outputs.

2. A program is not complete without sufficient testing to prove it's complete.

3. Global state and global impure functions makes testing more difficult. This
   can include singletons and system calls.

Any of these may be argued, but that will be left for other posts. Any of these
may be said of other types of programs as well, but that can also be left for
other posts.

## Components

Properties of components include:

1. *Creatable*: An instance of a component, given some defined set of
   parameters, can be created independently of any other instance of that or any
   other component.

2. *Composable*: A component may be used as a parameter of another component's
   instantiation. This would make it a child component of the one being
   instantiated (i.e. the parent).

3. *Abstract*: A component is an interface consisting of one or more methods.
   Being an interface, a component may have one or more implementations, but
   generally will have a primary implementation, which is used during a
   program's runtime, and secondary "mock" implementations, which are only used
   when testing other components.

4. *Isolated*: A component may not use mutable global variables (i.e.
   singletons) or impure global functions (e.g. system calls). It may only use
   constants and variables/components given to it during instantiation.

5. *Ephemeral*: A component may have a specific method used to clean up all
   resources that it's holding (e.g. network connections, file handles,
   language-specific lightweight threads, etc).

   5a. This cleanup method should _not_ clean up any child components given as
   instantiation parameters.

   5b. This cleanup method should not return until the component's cleanup is
   complete.

Components are composed together to create programs. This is done by passing
components as parameters to other components during instantiation. The `main`
process of the program is responsible for instantiating and composing most, if
not all, components in the program.

A component oriented program is one which primarily, if not entirely, uses
components for its functionality. Components generally have the quality of being
able to interact with code written in other patterns without any toes being
stepped on.

## Example

Let's start with an example: suppose a program is desired which accepts a string
over stdin, hashes it, then writes the string to a file whose name is the hash.

A naive implementation of this program in go might look like:

```go
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

func hashFileWriter() error {
	h := sha1.New()
	r := io.TeeReader(os.Stdin, h)
	body, _ := ioutil.ReadAll(r)
	fileName := hex.EncodeToString(h.Sum(nil))

	if err := ioutil.WriteFile(fileName, body, 0644); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := hashFileWriter(); err != nil {
		panic(err) // consider the error handled
	}
}
```

Notice that there's not a clear separation here between different components;
`hashFileWriter` _might_ be considered a one method component, except that it
breaks component property 4, which says that a component may not use mutable
global variables (`os.Stdin`) or impure global functions (`ioutil.WriteFile`).

Notice also that testing the program would require integration tests, and could
not be unit tested (because there are no units, i.e. components). For a trivial
program like this one writing unit and integration tests would be redundant, but
for larger programs it may not be. Unit tests are important because they are
fast to run, (usually) easy to formulate, and yield consistent results.

This program could instead be written as being composed of three components:

* `stdin`, a construct given by the runtime which outputs a stream of bytes.

* `disk`, accepts a file name and file contents as input, writes the file
  contents to a file of the given name, and potentially returns an error back.

* `hashFileWriter`, reads a stream of bytes off a `stdin`, collects the stream
  into a string, hashes that string to generate a file name, and uses `disk` to
  create a corresponding file with the string as its contents. If `disk` returns
  an error then `hashFileWriter` returns that error.

Sprucing up our previous example to use these more clearly defined components
might look like:

```go
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Disk defines the methods of the disk component.
type Disk interface {
	WriteFile(fileName string, fileContents []byte) error
}

// disk is the primary implementation of Disk. It implements the methods of
// Disk (WriteFile) by performing actual system calls.
type disk struct{}

func NewDisk() Disk { return disk{} }

func (disk) WriteFile(fileName string, fileContents []byte) error {
	return ioutil.WriteFile(fileName, fileContents, 0644)
}

func hashFileWriter(stdin io.Reader, disk Disk) error {
	h := sha1.New()
	r := io.TeeReader(stdin, h)
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	fileName := hex.EncodeToString(h.Sum(nil))

	if err := disk.WriteFile(fileName, body); err != nil {
		return fmt.Errorf("writing to file %q: %w", fileName, err)
	}
	return nil
}

func main() {
	if err := hashFileWriter(os.Stdin, NewDisk()); err != nil {
		panic(err) // consider the error handled
	}
}
```

`hashFileWriter` no longer directly uses `os.Stdin` and `ioutil.WriteFile`, but
instead takes in components wrapping them; `io.Reader` is a built-in interface
which `os.Stdin` inherently implements, and `Disk` is a simple interface defined
just for this program.

At first glance this would seem to have doubled the line-count for very little
gain. This is because we have not yet written tests.

## Testing

As has already been firmly established, testing is important.

In the second form of the program we can test the core-functionality of the
`hashFileWriter` component without resorting to using the actual `stdin` and
`disk` components. Instead we use mocks of those components. A mock component
implements the same input/outputs that the "real" component does, but in a way
which makes testing a particular component possible without reaching outside the
process. These are unit tests.

Tests for the latest form of the program might look like this:

```go
package main

import (
	"strings"
	"testing"
)

// mockDisk implements the Disk interface. When WriteFile is called mockDisk
// will pretend to write the file, but instead will simply store what arguments
// WriteFile was called with.
type mockDisk struct {
	fileName     string
	fileContents []byte
}

func (d *mockDisk) WriteFile(fileName string, fileContents []byte) error {
	d.fileName = fileName
	d.fileContents = fileContents
	return nil
}

func TestHashFileWriter(t *testing.T) {
	type test struct {
		in          string
		expFileName string
		// expFileContents can be inferred from in
	}

	tests := []test{
		{
			in:          "",
			expFileName: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			in:          "hello",
			expFileName: "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d",
		},
		{
			in:          "hello\nworld", // make sure newlines don't break things
			expFileName: "7db827c10afc1719863502cf95397731b23b8bae",
		},
	}

	for _, test := range tests {
		// stdin is mocked via a strings.Reader, which outputs the string it was
		// initialized with as a stream of bytes.
		in := strings.NewReader(test.in)

		// Disk is mocked by mockDisk, go figure.
		disk := new(mockDisk)

		if err := hashFileWriter(in, disk); err != nil {
			t.Errorf("in:%q got err:%v", test.in, err)
		} else if string(disk.fileContents) != test.in {
			t.Errorf("in:%q got contents:%q", test.in, disk.fileContents)
		} else if string(disk.fileName) != test.expFileName {
			t.Errorf("in:%q got fileName:%q", test.in, disk.fileName)
		}
	}
}
```

Notice that these tests do not _completely_ cover the desired functionality of
the program: if `disk` returns an error that error should be returned from
`hashFileWriter`. Whether or not this must be tested as well, and indeed the
pedantry level of tests overall, is a matter of taste. I believe these to be
sufficient.

## Configuration

Practically all programs require some level of runtime configuration. This may
take the form of command-line arguments, environment variables, configuration
files, etc. Almost all configuration methods will require some system call, and
so any component accessing configuration directly would likely break component
property 4.

Instead each component should take in whatever configuration parameters it needs
during instantiation, and let `main` handle collecting all configuration from
outside of the process and instantiating the components appropriately.

Let's take our previous program, but add in two new desired behaviors: first,
there should be a command-line parameter which allows for specifying the string
on the command-line, rather than reading from stdin, and second, there should be
a command-line parameter declaring which directory to write files into. The new
implementation looks like:

```
package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Disk defines the methods of the disk component.
type Disk interface {
	WriteFile(fileName string, fileContents []byte) error
}

// disk is the concrete implementation of Disk. It implements the methods of
// Disk (WriteFile) by performing actual OS calls.
type disk struct {
	dir string
}

func NewDisk(dir string) Disk { return disk{dir: dir} }

func (d disk) WriteFile(fileName string, fileContents []byte) error {
	fileName = filepath.Join(d.dir, fileName)
	return ioutil.WriteFile(fileName, fileContents, 0644)
}

func hashFileWriter(in io.Reader, disk Disk) error {
	h := sha1.New()
	r := io.TeeReader(in, h)
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	fileName := hex.EncodeToString(h.Sum(nil))

	if err := disk.WriteFile(fileName, body); err != nil {
		return fmt.Errorf("writing to file %q: %w", fileName, err)
	}
	return nil
}

func main() {
	str := flag.String("str", "", "If set, hash and write this string instead of stdin")
	dir := flag.String("dir", ".", "Directory which files should be written to")
	flag.Parse()

	var in io.Reader
	if *str == "" {
		in = os.Stdin
	} else {
		in = strings.NewReader(*str)
	}

	disk := NewDisk(*dir)

	if err := hashFileWriter(in, disk); err != nil {
		panic(err) // consider the error handled
	}
}
```

Very little has changed, and in fact `hashFileWriter` was not touched at all,
meaning all unit tests remained valid. 

## Setup/Runtime/Cleanup

A program can be split into three stages: setup, runtime, and cleanup. Setup
is the stage during which internal state is assembled in order to make runtime
possible. Runtime is the stage during which a program's actual function is being
performed. Cleanup is the stage during which runtime stop and internal state is
disassembled.

A graceful (i.e. reliably correct) setup is quite natural to accomplish, but
unfortunately a graceful cleanup is not a programmer's first concern, and
frequently is not a concern at all. However, when building reliable and correct
programs, a graceful cleanup is as important as a graceful setup and runtime. A
program is still running while it is being cleaned up, and it's possibly even
acting on the outside world still. Shouldn't it behave correctly during that
time?

Achieving a graceful setup and cleanup with components is quite simple:

During setup a single-threaded process (usually `main`) will construct the
"leaf" components (those which have no child components of their own) first,
then the components which take those leaves as parameters, then the components
which take _those_ as parameters, and so on, until all are constructed. The
components end up assembled into a directed acyclic graph.

At this point the program will begin runtime.

Once runtime is over and it is time for the program to exit it's only necessary
to call each component's cleanup method(s) in the reverse of the order the
components were instantiated in. A component's cleanup method should not be
called until all of its parent components have been cleaned up.

Inherent to the pattern is the fact that each component will certainly be
cleaned up before any of its child components, since its child components must
have been instantiated first and a component will not clean up child components
given as parameters (as-per component property 5a).

With go this pattern can be achieved easily using `defer`, but writing it out
manually is not so hard, as in this toy example:

```
package main

import (
	"fmt"
	"time"
)

// sleeper is a component which prints its children and sleeps when it's time to
// cleanup.
type sleeper struct {
	children []*sleeper
	toSleep  time.Duration

	// The builtin time.Sleep is an impure global function, a component can't
	// use it, so the component must be instantiated with it as a parameter.
	sleep func(time.Duration)

	// likewise os.Stdout is a global singleton, and so must also be a
	parameter.
	stdout io.Writer
}

func (s *sleeper) print() {
	fmt.Fprintf(s.stdout, "I will sleep for %v\n", s.toSleep)
	for _, child := range s.children {
		child.print()
	}
}

func (s *sleeper) cleanup() {
	s.sleep(s.toSleep)
	fmt.Fprintf(s.stdout, "I slept for %v\n", s.toSleep)
}

func main() {

	// Within main we make a helper function to easily construct sleepers. for a
	// toy like this it's not worth the effort of giving sleeper a real
	// initialization function.
	newSleeper := func(toSleep time.Duration, children ...*sleeper) *sleeper {
		return &sleeper{
			children: children,
			toSleep:  toSleep,
			sleep:    time.Sleep,
			stdout:   os.Stdout,
		}
	}

	aa := newSleeper(250 * time.Millisecond)
	defer aa.cleanup()

	ab := newSleeper(250 * time.Millisecond)
	defer ab.cleanup()

	// A's children are AA and AB
	a := newSleeper(500*time.Millisecond, aa, ab)
	defer a.cleanup()

	b := newSleeper(750 * time.Millisecond)
	defer b.cleanup()

	// root's children are A and B
	root := newSleeper(1*time.Second, a, b)
	defer root.cleanup()

	// All components are now instantiated and runtime begins.
	root.print()
    // ... and just like that, runtime ends.
	fmt.Println("--- Alright, fun is over, time for bed ---")

	// Now to clean up, cleanup methods are called in the reverse order of the
	// component's instantiation.
	root.cleanup()
	b.cleanup()
	a.cleanup()
	ab.cleanup()
	aa.cleanup()

	// Expected output is:
	//
	// I will sleep for 1s
	// I will sleep for 500ms
	// I will sleep for 250ms
	// I will sleep for 250ms
	// I will sleep for 750ms
	// --- Alright, fun is over, time for bed ---
	// I slept for 1s
	// I slept for 750ms
	// I slept for 500ms
	// I slept for 250ms
	// I slept for 250ms
}
```

## Criticisms

In lieu of a FAQ I will attempt to premeditate criticisms of the component
oriented pattern laid out in this post:

*This seems like a lot of extra work.*

Building reliable programs is a lot of work, just as building reliable-anything
is a lot of work. Many of us work in an industry which likes to balance
reliability (sometimes referred to by the more specious "quality") with
maleability and deliverability, which naturally leads to skepticism of any
suggestions which require more time spent on reliability. This is not
necessarily a bad thing, it's just how the industry functions.

All that said, a pattern need not be followed perfectly to be worthwhile, and
the amount of extra work incurred by it can be decided based on practical
considerations. I merely maintain that when it comes time to revisit some
existing code, either to fix or augment it, that the job will be notably easier
if the code _mostly_ follows this pattern.

*My language makes this difficult.*

I don't know of any language which makes this pattern particularly easy, so
unfortunately we're all in the same boat to some extent (though I recognize that
some languages, or their ecosystems, make it more difficult than others). It
seems to me that this pattern shouldn't be unbearably difficult for anyone to
implement in any language either, however, as the only language feature needed
is abstract typing.

It would be nice to one day see a language which explicitly supported this
pattern by baking the component properties in as compiler checked rules.

*This will result in over-abstraction.*

Abstraction is a necessary tool in a programmer's toolkit, there is simply no
way around it. The only questions are "how much?" and "where?".

The use of this pattern does not effect how those questions are answered, but
instead aims to more clearly delineate the relationships and interactions
between the different abstracted types once they've been established using other
methods. Over-abstraction is the fault of the programmer, not the language or
pattern or framework.

*The acronymn is CoP.*

Why do you think I've just been ackwardly using "this pattern" instead of the
acronymn for the whole post? Better names are welcome.

## Conclusion

The component oriented pattern helps make our code more reliable with only a
small amount of extra effort incurred. In fact most of the pattern has to do
establishing sensible abstractions around global functionality and remembering
certain idioms for how those abstractions should be composed together, something
most of us do to some extent already anyway.

While beneficial in many ways, component oriented programming is merely a tool
which can be applied in many cases. It is certain that there are cases where it
is not the right tool for the job. I've found these cases to be
few-and-far-between, however. It's a solid pattern that I've gotten good use out
of, and hopefully you'll find it, or some parts of it, to be useful as well.
