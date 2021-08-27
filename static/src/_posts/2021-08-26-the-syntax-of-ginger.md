---
title: >-
    The Syntax of Ginger
description: >-
    Oh man, this got real fun real quick.
series: ginger
tags: tech
---

Finally I have a syntax for ginger that I'm happy with. This has actually been a
huge roadblock for me up till this point. There's a bit of a chicken-and-the-egg
problem with the syntax: without pinning down the structures underlying the
syntax it's difficult to develop one, but without an idea of syntax it's
difficult to know what structures will be ergonomic to use.

I've been focusing on the structures so far, and have only now pinned down the
syntax. Let's see what it looks like.

## Preface: Conditionals

I've so far written [two][cond1] [posts][cond2] regarding conditionals in
ginger. After more reflection, I think I'm going to stick with my _original_
gut, which was to only have value and tuple vertices (no forks), and to use a
function which accepts both a boolean and two input edges: the first being the
one to take if the boolean is true, and the second being the one to take if it's
false.

Aka, the very first proposal in the [first post][cond1]. It's hard to justify
up-front, but I think once you see it in action with a clean syntax you'll agree
it just kind of works.

[cond1]: {% post_url 2021-03-01-conditionals-in-ginger %}
[cond2]: {% post_url 2021-03-04-conditionals-in-ginger-errata %}

## Designing a Syntax

Ginger is a bit of a strange language. It uses strange datastructures in strange
ways. But approaching the building of a syntax for any language is actually
straightforward: you're designing a serialization protocol.

To pull back a bit, consider a list of words. How would you encode this list in
order to write it to a file? To answer this, let's flip the question: how would
you design a sequence of characters (ie the contents of the file) such that the
reader could reconstruct the list?

Well, constructing the list from a sequence of characters requires being able to
construct it _at all_, so in what ways is the list constructed? For this list,
let's say there's only an append operation, which accepts a list and a value to
append to it, and returns the result.

If we say that append is encoded via wrapping parenthesis around its two
arguments, and that `()` encodes the empty list, then we get a syntax like...

```
(((() foo) bar) baz)
```

...which, in this instance, decodes to a list containing the words, "foo", "bar",
and "baz", in that order.

It's not a pretty syntax, but it demonstrates the method. If you know how the
datastructure is constructed via code, you know what capabilities the syntax must
have and how it needs to fit together.

## gg

All of this amounted to me needing to implement the ginger graph in some other
language, in order to see what features the syntax must have.

A few years ago I had begun an implementation of a graph datastructure in go, to
use as the base (or at least a reference) for ginger. I had called this
implementation `gg` (ginger graph), with the intention that this would also be
the file extension used to hold ginger code (how clever).

The basic qualities I wanted in a graph datastructure for ginger were, and still
are:

* Immutability, ie all operations which modify the structure should return a
  copy, leaving the original intact.

* Support for tuples.

* The property that it should be impossible to construct an invalid graph. An
  invalid graph might be, for example, one where there is a single node with no
  edges.

* Well tested, and reasonably performant.

Coming back to all this after a few years I had expected to have a graph
datastructure implemented, possibly with immutability, but lacking in tuples and
tests. As it turns out I completely underestimated my past self, because as far
as I can tell I had already finished the damn thing, tuples, tests and all.

It looks like that's the point where I stopped, probably for being unsure about
some other aspect of the language, and my motivation fell off. The fact that
I've come back to ginger, after all these years, and essentially rederived the
same language all over again, gives me a lot of confidence that I'm on the right
track (and a lot of respect for my past self for having done all this work!)

The basic API I came up with for building ginger graphs (ggs) looks like this:

```go
package gg

// OpenEdge represents an edge with a source value but no destination value,
// with an optional value on it. On its own an OpenEdge has no meaning, but is
// used as a building block for making Graphs.
type OpenEdge struct{ ... }

// TupleOut constructs an OpenEdge leading from a tuple, which is comprised of
// the given OpenEdges leading into it, with an optional edge value.
func TupleOut(ins []OpenEdge, edgeVal Value) OpenEdge

// ValueOut constructs an OpenEdge leading from a non-tuple value, with an
// optional edge value.
func ValueOut(val, edgeVal Value) OpenEdge

// ZeroGraph is an empty Graph, from which all Graphs are constructed via calls
// to AddValueIn.
var ZeroGraph = &Graph{ ... }

// Graph is an immutable graph structure, formed from a collection of edges
// between values and tuples.
type Graph struct{ ... }

// AddValueIn returns a new Graph which is a copy of the original, with the
// addition of a new edge. The new edge's source and edge value come from the
// given OpenEdge, and the edge's destination value is the given value.
func (g *Graph) AddValueIn(oe OpenEdge, val Value) *Graph
```

The actual API is larger than this, and includes methods to remove edges,
iterate over edges and values, and perform unions and disjoins of ggs. But the
above are the elements which are required only for _making_ ggs, which is all
that a syntax is concerned with.

As a demonstration, here is how one would construct the `min` operation, which
takes two numbers and returns the smaller, using the `gg` package:

```go
// a, b, in, out, if, etc.. are Values which represent the respective symbol.

// a is the result of passing in to the 0 operation, ie a is the 0th element of
// the in tuple.
min := gg.ZeroGraph.AddValueIn(gg.ValueOut(in, 0), a)

// b is the 1st element of the in tuple
min = min.AddValueIn(gg.ValueOut(in, 1), b)

// out is the result of an if which compares a and b together, and which returns
// the lesser.
min = min.AddValueIn(out, gg.TupleOut([]gg.OpenEdge{
    gg.TupleOut([]gg.OpenEdge{a, b}, lt),
    a,
    b,
}, if)
```

And here's a demonstration of how this `min` would be used:

```go
// out is the result of passing 1 and 5 to the min operation.
gg.ZeroGraph.AddValueIn(gg.TupleOut([]gg.OpenEdge{1, 5}, min), out)
```

## Make it Nice

_Technically_ we're done. We have an implementation of the language's underlying
structure, and a syntax which encodes it (ie the ugly ass go syntax above). But
obviously I'm not proposing anyone actually use that.

Another thing I found when digging around in the old ginger repo was a text
file, tucked away in a directory called "sandbox", which had a primitive syntax
which _almost_ worked. I won't copy it here, but you can find it if you care to.
But with that as a foundation I came up with a crude, rough draft spec, which
maps the go syntax to the new syntax.

```
ValueOut(val, edgeVal)    : -edgeVal-val
ValueOut(val, null)       : -val
TupleOut([]val, edgeVal)  : -edgeVal-(val, ...)
TupleOut([]val, null)     : -(val, ...)
Graph(openEdge->val, ...) : { val openEdge, ... }
```

A couple things to note about this spec:

* `null` is used to indicate absence of value on an edge. The details of `null`
  are yet to be worked out, but we can use this placeholder for now.

* `Graph` is cheating a bit. In the original `gg` implementation a Graph gains
  its OpenEdge/Value pairs via successive calls to `AddValueIn`. However, such a
  pattern doesn't translate well to text, and since we're dealing purely with
  constructing an entire Graph at once we can instead have our Graph syntax
  declare all OpenEdge/Value pairs at once.

* It's backwards! Eg where the go syntax does `ValueOut(val, edgeVal)`, the
  proposed spec puts the values in the opposite order: `-edgeVal-val`. The
  former results in code which is read from input to output, while the latter
  results in the opposite: output to input.

  This was a tip I picked up from the old text file I found, and the result is
  code which is more familiar to an existing programmer. I _think_ (but am
  not sure) that it's also more in line with how programming is done mentally,
  ie we start with a result and work backwards to figure out what it takes to
  get there.

  It's possible, though, that I'm wrong, so at this end of this post I'm going
  to put some examples of the same code both "forwards" and "backwards" and see
  how I feel about it.

With all that said, let's see it in action! Here's `min` implemented in our shiny new syntax:

```
min -{
    a -0-in,
    b -1-in,
    out -if-(
        -lt-(-a,-b),
        -a,
        -b
    )
}
```

and then here's it being used:

```
out -min-(-1,-5)
```

## Make it _Nicer_

The most striking feature of this rough draft spec is all the prefix dashes,
such as in the `-min-(-1,-5)` statement. These dashes were included as they make
sense in the context of what the intended human interpretation of the structure
is: two values, `1`, and `5`, are being _piped into_ the two slots of a 2-tuple,
and that 2-tuple is being _piped into_ the `min` operation, the output of which
is being _piped into_ something `out`.

The "piping into" is what the dash represents, which is why the top level values
in the graph, `a`, `b`, and `out`, don't have a preceding dash; they are the
ultimate destinations of the pipes leading to them. But these pipes are
ultimately ugly, and also introduce odd questions like "how do we represent
-1?", so they need to go.

So I've made a second draft, which is only a few changes away from the rough,
but oh man do those changes make a world of difference. Here's the cleaned up
spec:

```
ValueOut(val, edgeVal)    : edgeVal(val)
ValueOut(val, null)       : val
TupleOut([]val, edgeVal)  : edgeVal(val, ...)
TupleOut([]val, null)     : (val, ...)
Graph(openEdge->val, ...) : { val = openEdge, ... }
```

The dashes were simply removed, and the `edgeVal` and `val` concatted together.
For `ValueOut(val, edgeVal)` wrapping parenthesis were put around `val`, to
delineate it and `edgeVal`. This conflicts with the syntax for `TupleOut([]val,
edgeVal)`, but that conflict is easy to remedy: when parenthesis wrap only a
single `val` then that is a `ValueOut`, otherwise it's a `TupleOut`.

Another change is to add an `=` between the `val` and `openEdge` in the `Graph`
constructor. This is a purely aesthetic change, but as you'll see it works well.

So let's see it! `min` implemented with this cleaned up syntax:

```
min = {
    a = 0(in),
    b = 1(in),
    out = if(
        lt(a,b),
        a,
        b
    )
}
```

And then its use:

```
min(1,5)
```

Well well well, look what we have here: a conventional programming language
syntax! `{`/`}` wrap a scope, and `(`/`)` wrap function arguments and
(optionally) single values. It's a lot clearer now that `0` and `1` are being
used as operations themselves when instantiating `a` and `b`, and `if` is much
more readable.

I was extremely surprised at how well this actually worked out. Despite having
drastically different underpinnings than most languages it ends up looking both
familiar and obvious. How cool!

## Examples Examples Examples

Here's a collection of example programs written in this new syntax. The base
structure of these are borrowed from previous posts, I'm merely translating that
structure into a new form:

```
// decr outputs one less than the input.
decr = { out = add(in, -1) }

// fib accepts a number i, and outputs the ith fibonacci number.
fib = {

    inner = {
        n = 0(in),
        a = 1(in),
        b = 2(in),

        out = if(zero?(n),
            a,
            inner(decr(n), b, add(a,b))
        )

    },

    out = inner(in, 0, 1)
}

// map accepts a sequence and a function, and returns a sequence consisting of
// the result of applying the function to each of the elements in the given
// sequence.
map = {
    inner = {
        mapped-seq = 0(in),
        orig-seq = 1(in),
        op = 2(in),

        i = len(mapped-seq),

        // graphs provide an inherent laziness to the language. Just because
        // next-el is _defined_ here doesn't mean it's evaluated here at runtime.
        // In reality it will only be evaluated if/when evaluating out requires
        // evaluating next-el.
        next-el = op(i(orig-seq)),
        next-mapped-seq = append(mapped-seq, next-el),

        out = if(
            eq(len(mapped-seq), len(orig-seq)),
            mapped-seq,
            inner(next-mapped-seq, orig-seq, op)
        )
    }

    // zero-seq returns an empty sequence
    out = inner(zero-seq(), 0(in), 1(in))
}
```

## Selpmexa Selpmexa Selpmexa

Our syntax encodes a graph, and a graph doesn't really care if the syntax was
encoded in an input-to-output vs an output-to-input direction. So, as promised,
here's all the above examples, but "backwards":

```
// min returns the lesser of the two numbers it is given
{
    (in)0 = a,
    (in)1 = b,

    (
        (a,b)lt,
        a,
        b
    )if = out

} = min

// decr outputs one less than the input.
{ (in, -1)add = out } = decr

// fib accepts a number i, and outputs the ith fibonacci number.
{
    {
        (in)0 = n,
        (in)1 = a,
        (in)2 = b,

        (
            (n)zero?
            a,
            ((n)decr, b, (a,b)add)inner
        )if = out

    } = inner,

    (in, 0, 1)inner = out

} = fib

// map accepts a sequence and a function, and returns a sequence consisting of
// the result of applying the function to each of the elements in the given
// sequence.
{
    {
        (in)0 = mapped-seq,
        (in)1 = orig-seq,
        (in)2 = op,

        (mapped-seq)len = i,

        ((orig-seq)i)op = next-el,
        (mapped-seq, next-el)append = next-mapped-seq,

        (
            ((mapped-seq)len, (orig-seq)len)eq,
            mapped-seq,
            (next-mapped-seq, orig-seq, op)inner
        )if = out

    } = inner,

    (()zero-seq, (in)0, (in)1)inner = out
} = map
```

Do these make you itchy? They kind of make me itchy. But... parts of them also
appeal to me.

The obvious reason why these feel wrong to me is the placement of `if`:

```
(
    (a,b)lt,
    a,
    b
)if = out
```

The tuple which is being passed to `if` here is confusing unless you already
know that it's going to be passed to `if`. But on your first readthrough you
won't know that till you get to the end, so you'll be in the dark until then.
For more complex programs I'm sure this problem will compound.

On the other hand, pretty much everything else looks _better_, imo. For example:

```
// copied and slightly modified from the original to make it even more complex

(mapped-seq, ((orig-seq)i)op)append = next-mapped-seq
```

Something like this reads very clearly to me, and requires a lot less mental
backtracking to comprehend. The main difficulty I have is tracking the
parenthesis, but the overall "flow" of data and the order of events is plain to
read.

## Next Steps

The syntax here is not done yet, not by a long shot. If my record with past
posts about ginger (wherein I've "decided" on something and then completely
backtracked in later posts every single time) is any indication then this syntax
won't even look remotely familiar in a very short while. But it's a great
starting point, I think, and raises a lot of good questions.

* Can I make parenthesis chains, a la the last example, more palatable in some
  way?

* Should I go with the "backwards" syntax afterall? In a functional style of
  programming `if` statements _should_ be in the minority, and so the syntax
  which better represents the flow of data in that style might be the way.

* Destructuring of tuples seems to be wanted, as evidenced by all the `a =
  0(in)` lines. Should this be reflected in the structure or solely be
  syntactical sugar?

* Should the commas be replaced with any whitespace (and make commas count as
  whitespace, as clojure has done)? If this is possible then I think they should
  be, but I won't know for sure until I begin implementing the parser.

And, surely, many more! I've felt a bit lost with ginger for a _long_ time, but
seeing a real, usable syntax emerge has really invigorated me, and I'll be
tackling it again in earnest soon (fingers crossed).
