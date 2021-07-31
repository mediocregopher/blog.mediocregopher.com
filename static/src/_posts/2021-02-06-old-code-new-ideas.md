---
title: >-
    Old Code, New Ideas
description: >-
    Looking back at my old code with bemusement and horror.
tags: tech
---

About 3 years ago I put a lot of effort into a set of golang packages called
[mediocre-go-lib][mediocre-go-lib]. The idea was to create a framework around
the ideas I had laid out in [this blog post][program-structure] around the
structure and composability of programs. What I found in using the framework was
that it was quite bulky, not fully thought out, and ultimately difficult for
anyone but me to use. So.... a typical framework then.

My ideas about program structure haven't changed a ton since then, but my ideas
around the patterns which enable that structure have simplified dramatically
(see [my more recent post][component-oriented] for more on that). So in that
spirit I've decided to cut a `v2` branch of `mediocre-go-lib` and start trimming
the fat.

This is going to be an exercise both in deleting old code (very fun) and
re-examining old code which I used to think was good but now know is bad (even
more fun), and I've been looking forward to it for some time.

[mediocre-go-lib]: https://github.com/mediocregopher/mediocre-go-lib
[program-structure]: {% post_url 2019-08-02-program-structure-and-composability %}
[component-oriented]: {% post_url 2020-11-16-component-oriented-programming %}

## mcmp, mctx

The two foundational pieces of `mediocre-go-lib` are the `mcmp` and `mctx`
packages. `mcmp` primarily deals with its [mcmp.Component][component] type,
which is a key/value store which can be used by other packages to store and
retrieve component-level information. Each `mcmp.Component` exists as a node in
a tree of `mcmp.Component`s, and these form the structure of a program.
`mcmp.Component` is able to provide information about its place in that tree as
well (i.e. its path, parents, children, etc...).

If this sounds cumbersome and of questionable utility that's because it is. It's
also not even correct, because a component in a program exists in a DAG, not a
tree. Moreover, each component can keep track of whatever data it needs for
itself using typed fields on a struct. Pretty much all other packages in
`mediocre-go-lib` depend on `mcmp` to function, but they don't _need_ to, I just
designed it that way.

So my plan of attack is going to be to delete `mcmp` completely, and repair all
the other packages.

The other foundational piece of `mediocre-go-lib` is [mctx][mctx]. Where `mcmp`
dealt with arbitrary key/value storage on the component level, `mctx` deals with
it on the contextual level, where each go-routine (i.e. thread) corresponds to a
`context.Context`. The primary function of `mctx` is this one:

```go
// Annotate takes in one or more key/value pairs (kvs' length must be even) and
// returns a Context carrying them.
func Annotate(ctx context.Context, kvs ...interface{}) context.Context
```

I'm inclined to keep this around for now because it will be useful for logging,
but there's one change I'd like to make to it. In its current form the value of
every key/value pair must already exist before being used to annotate the
`context.Context`, but this can be cumbersome in cases where the data you'd want
to annotate is quite hefty to generate but also not necessarily going to be
used. I'd like to have the option to make annotating occur lazily.  For this I
add an `Annotator` interface and a `WithAnnotator` function which takes it as an
argument, as well as some internal refactoring to make it all work right:

```go
// Annotations is a set of key/value pairs representing a set of annotations. It
// implements the Annotator interface along with other useful post-processing
// methods.
type Annotations map[interface{}]interface{}

// Annotator is a type which can add annotation data to an existing set of
// annotations. The Annotate method should be expected to be called in a
// non-thread-safe manner.
type Annotator interface {
	Annotate(Annotations)
}

// WithAnnotator takes in an Annotator and returns a Context which will produce
// that Annotator's annotations when the Annotations function is called. The
// Annotator will be not be evaluated until the first call to Annotations.
func WithAnnotator(ctx context.Context, annotator Annotator) context.Context
```

`Annotator` is designed like it is for two reasons. The more obvious design,
where the method has no arguments and returns a map, would cause a memory
allocation on every invocation, which could be a drag for long chains of
contexts whose annotations are being evaluated frequently. The obvious design
also leaves open questions about whether the returned map can be modified by
whoever receives it. The design given here dodges these problems without any
obvious drawbacks.

The original implementation also had this unnecessary `Annotation` type:

```go
// Annotation describes the annotation of a key/value pair made on a Context via
// the Annotate call.
type Annotation struct {
       Key, Value interface{}
}
```

I don't know why this was ever needed, as an `Annotation` was never passed into
nor returned from any function. It was part of the type `AnnotationSet`, but
that could easily be refactored into a `map[interface{}]interface{}` instead. So
I factored `Annotation` out completely.

[component]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mcmp#Component
[mctx]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mctx

## mcfg, mrun

The next package to tackle is [mcfg][mcfg], which deals with configuration via
command line arguments and environment variables. The package is set up to use
the old `mcmp.Component` type such that each component could declare its own
configuration parameters in the global configuration. In this way the
configuration would have a hierarchy of its own which matches the component
tree.

Given that I now think `mcmp.Component` isn't the right course of action it
would be the natural step to take that aspect out of `mcfg`, leaving only a
basic command-line and environment variable parser. There are many other basic
parsers of this sort out there, including [one][flagconfig] or [two][lever] I
wrote myself, and frankly I don't think the world needs another. So `mcfg` is
going away.

The [mrun][mrun] package is the corresponding package to `mcfg`; where `mcfg`
dealt with configuration of components `mrun` deals with the initialization and
shutdown of those same components. Like `mcfg`, `mrun` relies heavily on
`mcmp.Component`, and doesn't really have any function with that type gone. So
`mrun` is a gonner too.

[mcfg]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mcfg
[mrun]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mrun
[flagconfig]: https://github.com/mediocregopher/flagconfig
[lever]: https://github.com/mediocregopher/lever

## mlog

The [mlog][mlog] package is primarily concerned with, as you might guess,
logging.  While there are many useful logging packages out there none of them
integrate with `mctx`'s annotations, so it is useful to have a custom logging
package here. `mlog` also has the nice property of not being extremely coupled
to `mcmp.Component` like other packages; it's only necessary to delete a handful
of global functions which aren't a direct part of the `mlog.Logger` type in
order to free the package from that burden.

With that said, the `mlog.Logger` type could still use some work. It's primary
pattern looks like this:

```go
// Message describes a message to be logged.
type Message struct {
	Level
	Description string
	Contexts []context.Context
}

// Info logs an InfoLevel message.
func (l *Logger) Info(descr string, ctxs ...context.Context) {
	l.Log(mkMsg(InfoLevel, descr, ctxs...))
}
```

The idea was that if the user has multiple `Contexts` in hand, each one possibly
having some relevant annotations, all of those `Context`s' annotations could be
merged together for the log entry.

Looking back it seems to me that the only thing `mlog` should care about is the
annotations, and not _where_ those annotations came from. So the new pattern
looks like this:

```go
// Message describes a message to be logged.
type Message struct {
	Context context.Context
	Level
	Description string
	Annotators  []Annotators
}

// Info logs a LevelInfo message.
func (l *Logger) Info(ctx context.Context, descr string, annotators ...mctx.Annotator)
```

The annotations on the given `Context` will be included, and then any further
`Annotator`s can be added on. This will leave room for `merr` later.

There's some other warts in `mlog.Logger` that should be dealt with as well,
including some extraneous methods which were only used due to `mcmp.Component`,
some poorly named types, a message handler which didn't properly clean itself
up, and making `NewLogger` take in parameters with which it can be customized as
needed (previously it only allowed for a single configuration). I've also
extended `Message` to include a timestamp, a namespace field, and some other
useful information.

[mlog]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mlog

## Future Work

I've run out of time for today, but future work on this package includes:

* Updating [merr][merr] with support for `mctx.Annotations`.
* Auditing the [mnet][mnet], [mhttp][mhttp], and [mrpc][mrpc] packages to see if
  they contain anything worth keeping.
* Probably deleting the [m][m] package entirely; I don't even really remember
  what it does.
* Probably deleting the [mdb][mdb] package entirely; it only makes sense in the
  context of `mcmp.Component`.
* Making a difficult decision about [mtest][mtest]; I put a lot of work into it,
  but is it really any better than [testify][testify]?

[merr]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/merr
[mnet]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mnet
[mhttp]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mhttp
[mrpc]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mrpc
[m]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/m
[mdb]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mdb
[mtest]: https://pkg.go.dev/github.com/mediocregopher/mediocre-go-lib/mtest
[testify]: https://github.com/stretchr/testify
