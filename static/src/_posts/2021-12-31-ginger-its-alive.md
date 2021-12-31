---
title: >-
    Ginger: It's Alive!
description: >-
    The new best language for computing fibonacci numbers.
series: ginger
tags: tech
---

As a kind of Christmas present to myself I took a whole week off of work
specifically to dedicate myself to working on ginger.

My concrete goal was to be able to run a ginger program to compute any Nth
fibonacci number, a goal I chose because it would require the implementation of
conditionals, some kind of looping or recursion, and basic addition/subtraction.
In other words, it would require all the elements which comprise a Turing
complete language.

And you know what? I actually succeeded!

The implementation can be found [here][impl]. At this point ginger is an
interpreted language running in a golang-based VM. The dream is for it to be
self-hosted on LLVM (and other platforms after), but as an intermediate step to
that I decided on sticking to what I know (golang) rather than having to learn
two things at once.

In this post I'm going to describe the components of this VM at a high level,
show a quick demo of it working, and finally talk about the roadmap going
forward.

[impl]: https://github.com/mediocregopher/ginger/tree/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2

## Graph

The core package of the whole project is the [`graph`][graph] package. This
package implements a generic directed graph datastructure.

The generic part is worth noting; I was able to take advantage of go's new
generics which are currently [in beta][go118]. I'd read quite a bit on how the
generic system would work even before the beta was announced, so I was able to
hit the ground running and start using them without much issue.

Ginger's unique graph datastructure has been discussed in previous posts in this
series quite a bit, and this latest implementation doesn't deviate much at a
high level. Below are the most up-to-date core datatypes and functions which are
used to construct ginger graphs:

```go

// Value is any value which can be stored within a Graph. Values should be
// considered immutable, ie once used with the graph package their internal
// value does not change.
type Value interface {
	Equal(Value) bool
	String() string
}

// OpenEdge consists of the edge value (E) and source vertex value (V) of an
// edge in a Graph. When passed into the AddValueIn method a full edge is
// created. An OpenEdge can also be sourced from a tuple vertex, whose value is
// an ordered set of OpenEdges of this same type.
type OpenEdge[E, V Value] struct { ... }

// ValueOut creates a OpenEdge which, when used to construct a Graph, represents
// an edge (with edgeVal attached to it) coming from the vertex containing val.
func ValueOut[E, V Value](edgeVal E, val V) *OpenEdge[E, V]

// TupleOut creates an OpenEdge which, when used to construct a Graph,
// represents an edge (with edgeVal attached to it) coming from the vertex
// comprised of the given ordered-set of input edges.
func TupleOut[E, V Value](edgeVal E, ins ...*OpenEdge[E, V]) *OpenEdge[E, V]

// Graph is an immutable container of a set of vertices. The Graph keeps track
// of all Values which terminate an OpenEdge. E indicates the type of edge
// values, while V indicates the type of vertex values.
type Graph[E, V Value] struct { ... }

// AddValueIn takes a OpenEdge and connects it to the Value vertex containing
// val, returning the new Graph which reflects that connection.
func (*Graph[E, V]) AddValueIn(val V, oe *OpenEdge[E, V]) *Graph[E, V]

// ValueIns returns, if any, all OpenEdges which lead to the given Value in the
// Graph (ie, all those added via AddValueIn).
func (*Graph[E, V]) ValueIns(val Value) []*OpenEdge[E, V]

```

The current `Graph` implementation is _incredibly_ inefficient, it does a lot of
copying, looping, and equality checks which could be optimized out one day.
That's going to be a recurring theme of this post, as I had to perform a
balancing act between actually reaching my goal for the week while not incurring
too much tech debt for myself.

[graph]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/graph/graph.go
[go118]: https://go.dev/blog/go1.18beta1

### MapReduce

There's a final operation I implemented as part of the `graph` package:
[MapReduce][mapreduce]. It's a difficult operation to describe, but I'm going to
do my best in this section for those who are interested. If you don't understand
it, or don't care, just know that `MapReduce` is a generic tool for transforming
graphs.

For a description of `MapReduce` we need to present an example graph:

```
        +<--b---
        +       \
X <--a--+<--c----+<--f-- A
        +               /
        +      +<---g---
        +<--d--+
               +<---h---
                        \
Y <---------e----------- B
```

Plus signs indicate tuples, and lowercase letters are edge values while upper
case letters are vertex values. The pseudo-code to construct this graph in go
might look like:

```go
    g := new(Graph)

    fA := ValueOut("f", "A")

    g = g.AddValueIn(
        "X",
        TupleOut(
            "a",
            TupleOut("b", fA),
            TupleOut("c", fA),
            TupleOut(
                "d",
                ValueOut("g", "A"),
                ValueOut("h", "B"),
            ),
        ),
    )

    g = g.AddValueIn("e", "B")
```

As can be seen in the [code][mapreduce], `MapReduce`'s first argument is an
`OpenEdge`, _not_ a `Graph`. Fundamentally `MapReduce` is a reduction of the
_dependencies_ of a particular value into a new value; to reduce the
dependencies of multiple values at the same time would be equivalent to looping
over those values and calling `MapReduce` on each individually. Having
`MapReduce` only deal with one edge at a time is more flexible.

So let's focus on a particular `OpenEdge`, the one leading into `X` (returned by
`TupleOut("a", etc...)`. `MapReduce` is going to descend into this `OpenEdge`
recursively, in order to first find all value vertices (ie the leaf vertices,
those without any children of their own).

At this point `MapReduce` will use its second argument, the `mapVal` function,
which accepts a value of one type and returns a value of another type. This
function is called on each value from every value vertex encountered. In this
case both `A` and `B` are connectable from `X`, so `mapVal` will be called on
each _only once_. This is the case even though `A` is connected to multiple
times (once with an edge value of `f`, another with an edge value of `b`).
`mapVal` only gets called once per vertex, not per connection.

With all values mapped, `MapReduce` will begin reducing. For each edge leaving
each value vertex, the `reduceEdge` function is called. `reduceEdge` accepts as
arguments the edge value of the edge and the _mapped value_ (not the original
value) of the vertex, and returns a new value of the same type that `mapVal`
returned. Like `mapVal`, `reduceEdge` will only be called once per edge. In our
example, `<--f--A` is used twice (`b` and `c`), but `reduceEdge` will only be
called on it once.

With each value vertex edge having been reduced, `reduceEdge` is called again on
each edge leaving _those_ edges, which must be tuple edges. An array of the
values returned from the previous `reduceEdge` calls for each of the tuples'
input edges is used as the value argument in the next call. This is done until
the `OpenEdge` is fully reduced into a single value.

To flesh out our example, let's imagine a `mapVal` which returns the input
string repeated twice, and a `reduceEdge` which returns the input values joined
with the edge value, and then wrapped with the edge value (eg `reduceEdge(a, [B,
C]) -> aBaCa`).

Calling `MapReduce` on the edge leading into `X` will then give us the following
calls:

```
# Map the value vertices

mapVal(A) -> AA
mapVal(B) -> BB

# Reduce the value vertex edges

reduceEdge(f, [AA]) -> fAAf
reduceEdge(g, [AA]) -> gAAg
reduceEdge(h, [BB]) -> hBBh

# Reduce tuple vertex edges

reduceEdge(b, [fAAf]) -> bfAAfb
reduceEdge(c, [fAAf]) -> cfAAfc
reduceEdge(d, [gAAg, hBBh]) -> dgAAgdhBBhd

reduceEdge(a, [bfAAfb, cfAAfc, dgAAgdhBBhd]) -> abfAAfbacfAAfcadgAAgdhBBhda
```

Beautiful, exactly what we wanted.

`MapReduce` will prove extremely useful when it comes time for the VM to execute
the graph. It enables the VM to evaluate only the values which are needed to
produce an output, and to only evaluate each value once no matter how many times
it's used. `MapReduce` also takes care of the recursive traversal of the
`Graph`, which simplifies the VM code significantly.

[mapreduce]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/graph/graph.go#L338

## gg

With a generic graph implementation out of the way, it was then required to
define a specific implementation which could be parsed from a file and later
used for execution in the VM.

The file extension used for ginger code is `.gg`, as in "ginger graph" (of
course). The package name for decoding this file format is, therefore, also
called `gg`.

The core datatype for the `gg` package is the [`Value`][ggvalue], since the
`graph` package takes care of essentially everything else in the realm of graph
construction and manipulation. The type definition is:

```go
// Value represents a value which can be serialized by the gg text format.
type Value struct {

	// Only one of these fields may be set
	Name   *string
	Number *int64
	Graph  *Graph

	// Optional fields indicating the token which was used to construct this
	// Value, if any.
	LexerToken *LexerToken
}

type Graph = graph.Graph[Value, Value] // type alias for convenience
```

Note that it's currently only possible to describe three different types in a
`gg` file, and one of them is the `Graph`! These are the only ones needed to
implement a fibonacci function, so they're all I implemented.

The lexing/parsing of `gg` files is not super interesting, you can check out the
package code for more details. The only other thing worth noting is that, for
now, all statements are required to end with a `;`. I had originally wanted to
be less strict with this, and allow newlines and other tokens to indicate the
end of statements, but it was complicating the code and I wanted to move on.

Another small thing worth noting is that I decided to make each entire `.gg`
file implicitly define a graph. So you can imagine each file's contents wrapped
in curly braces.

With the `gg` package out of the way I was able to finally parse ginger
programs! The following is the actual, real-life implementation of the fibonacci
function (though at this point it didn't actually work, because the VM was still
not implemented:

```
out = {

    decr = { out = add < (in; -1;); };

    n = tupEl < (in; 0;);
    a = tupEl < (in; 1;);
    b = tupEl < (in; 2;);

    out = if < (
        isZero < n;
        a;
        recur < (
            decr < n;
            b;
            add < (a;b;);
        );
    );

} < (in; 0; 1;);
```

[ggvalue]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/gg/gg.go#L14

## VM

Finally, the meat of all this. If the `graph` and `gg` packages are the sturdy,
well constructed foundations of a tall building, then the `vm` package is the
extremely long, flimsy stick someone propped up vertically so they could say
they built a structure of impressive height.

In other words, it's very likely that the current iteration of the VM will not
be long for this world, and so I won't waste time describing it in super detail.

What I will say about it is that within the `vm` package I've defined a [new
`Value` type][vmvalue], which extends the one defined in `gg`. The necessity of
this was that there are types which cannot be represented syntactically in a
`.gg` file, but which _can_ be used as values within a program being run.

The first of these is the `Operation`, which is essentially a first-class
function. The VM will automatically interpret a graph as an `Operation` when it
is used as an edge value, as has been discussed in previous posts, but there are
also built-in operations (like `if` and `recur`) which cannot be represented as
datastructures, and so it was necessary to introduce a new in-memory type to
properly represent operations.

The second is the `Tuple` type. This may seem strange, as ginger graphs already
have a concept of a tuple. But the ginger graph tuple is a _vertex type_, not a
value type. The distinction is small, but important. Essentially the graph tuple
is a structural element which describes how to create a tuple value, but it is
not yet that value. So we need a new Value type to hold the tuple once it _has_
been created during runtime.

Another thing worth describing about the `vm` package, even though I think they
might change drastically, are [`Thunk`s][thunk]:

```go
// Thunk is returned from the performance of an Operation. When called it will
// return the result of that Operation having been called with the particular
// arguments which were passed in.
type Thunk func() (Value, error)
```

The term "thunk" is borrowed from Haskell, which I don't actually know so I'm
probably using it wrong, but anyway...

A thunk is essentially a value which has yet to be evaluated; the VM knows
exactly _how_ to evaluate it, but it hasn't done so yet. The primary reason for
their existence within ginger is to account for conditionals, ie the `if`
operation. The VM can't evaluate each of an `if`'s arguments all at once, it
must only evaluate the first argument (to obtain a boolean), and then based on
that evaluate the second or third argument.

This is where `graph.MapReduce` comes in. The VM uses `graph.MapReduce` to
reduce each edge in a graph to a `Thunk`, where the `Thunk`'s value is based on
the operation (the edge's value) and the inputs to the edge (which will
themselves be `Thunk`s). Because each `Thunk` represents a potential value, not
an actual one, the VM is able to completely parse the program to be executed
(using `graph.MapReduce`) while allowing conditionals to still work correctly.

[EvaluateEdge][evaledge] is where all that happens, if you're interested, but be
warned that the code is a hot mess right now and it's probably not worth
spending a ton of time understanding it as it will change a lot.

A final thing I'll mention is that the `recur` operation is, I think, broken. Or
probably more accurately, the entire VM is broken in a way which prevents
`recur` from working correctly. It _does_ produce the correct output, so I
haven't prioritized debugging it, but for any large number of iterations it
takes a very long time to run.

[vmvalue]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/vm/vm.go#L18
[thunk]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/vm/op.go#L11
[evaledge]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/vm/scope.go#L29

## Demo

Finally, to show it off! I put together a super stupid `eval` binary which takes
two arguments: a graph to be used as an operation, and a value to be used as an
argument to that operation. It doesn't even read the code from a file, you have
to `cat` it in.

The [README][readme] documents how to run the demo, so if you'd like to do so
then please clone the repo and give it a shot! It should look like this when you
do:

```
# go run ./cmd/eval/main.go "$(cat examples/fib.gg)" 8
21
```

You can put any number you like instead of `8`, but as mentioned, `recur` is
broken so it can take a while for larger numbers.

[readme]: https://github.com/mediocregopher/ginger/blob/ebf57591a8ac08da8a312855fc3a6d9c1ee6dcb2/README.md

## Next Steps

The following are all the things I'd like to address the next time I work on
ginger:

* `gg`

    * Allow for newlines (and `)` and `}`) to terminate statements, not just
      `;`.

    * Allow names to have punctuation characters in them (maybe?).

    * Don't read all tokens into memory prior to parsing.

* `vm`

    * Fix `recur`.

    * Implement tail call optimization.

* General

    * A bit of polish on the `eval` tool.

    * Expose graph creation, traversal, and transformation functions as
      builtins.

    * Create plan (if not actually implement it yet) for how code will be
      imported from one file to another. Namespacing in general will fall into
      this bucket.

    * Create plan (if not actually implement it yet) for how users can
      extend/replace the lexer/parser.

I don't know _when_ I'll get to work on these next, ginger will come back up in
my rotation of projects eventually. It could be a few months. In the meantime I
hope you're as excited about this progress as I am, and if you have any feedback
I'd love to hear it.

Thanks for reading!
