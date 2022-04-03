---
title: >-
    Ginger: A Small VM Update
description: >-
    It works gooder now.
tags: tech
series: ginger
---

During some recent traveling I had to be pulled away from cryptic-net work for a
while. Instead I managed to spend a few free hours, and the odd international
plane ride, to fix the ginger vm.

The problem, as it stood, was that it only functioned "correctly" in a very
accidental sense. I knew from the moment that I published it that it would get
mostly rewritten immediately.

And so here we are, with a rewritten vm and some new realizations.

## Operation

The `Operation` type was previously defined like so:

```
type Operation interface {
	Perform([]Thunk, Operation) (Thunk, error)
}
```

I'm not going to explain it, because it's both confusing and wrong.

One thing that is helpful in a refactor, especially in a strongly typed
language, is to tag certain interfaces as being axiomatic, and conforming the
rest of your changes around those. If those interfaces are simple enough to
apply broadly _and_ accurately describe desired behavior, they will help
pre-decide many difficult decisions you'd otherwise have to make.

So with that mind, I tagged `Operation` as being an axiomatic interface, given
that ginger is aiming to be a functional language (and I'm wondering if I should
just rename `Operation` to `Function`, while I'm at it). The new definition of
the interface is:

```
type Operation interface {
	Perform(Value) Value
}
```

`Operation` takes and argument and returns a result, it could not possibly be
boiled down any further. By holding `Operation` to this definition and making
decisions from there, it was pretty clear what the next point of attack was.

## If/Recur

The reason that `Operation` had previously been defined in such a fucked up way
was to support the `if` and `recur` `Operation`s, as if they weren't different
than any other `Operation`s. But truthfully they are different, as they are
actually control flow constructs, and so require capabilities that no other
`Operation` would be allowed to use anyway.

The new implementation reflects this. `if` and `recur` are now both handled
directly by the compiler, while global `Operation`s like `tupEl` are
implementations of the `Operation` interface.

## Compile Step

The previous iteration of the vm hadn't distinguished between a compile step and
a run step. In a way it did both at the same time, by abusing the `Thunk` type.
Separating the two steps, and ditching the `Thunk` type in the process, was the
next major change in the refactoring.

The compile step can be modeled as a function which takes a `Graph` and returns
an `Operation`, where the `Graph`'s `in` and `out` names correspond to the
`Operation`'s argument and return, respectively. The run step then reads an
input from the user, calls the compiled `Operation` with that input, and outputs
the result back to the user.

As an example, given the following program:

```
* six-or-more.gg

max = {
    a = tupEl < (in, 0)
    b = tupEl < (in, 1)
    out = if < (gt < (a, b), a, b)
}

out = max < (in, 6)
```

we want to compile an `Operation` which accepts a number and returns the greater
of that number and 6. I'm going to use anonymous go functions to demonstrate the
anatomy of the compiled `Operation`, as that's what's happening in the current
compiler anyway.

```
// After compilation, this function will be in-memory and usable as an
// Operation.

sixOrMore := func(in Value) Value {

    max := func(in Value) Value {

        a := tupEl(in, 0)
        b := tupEl(in, 1)

        if a > b {
            return a
        }

        return b
    }

    return max(in, 6)
}
```

Or at least, this is what I tried for _initially_. What I found was that it was
easier, in the context of how `graph.MapReduce` works, to make even the leaf
values, e.g. `in`, `0`, `1`, and `6`, map to `Operations` as well. `in` is
replaced with an anonymous function which returns its argument, and the numbers
are replaced with anonymous functions which ignore their argument and always
return their respective number.

So the compiled `Operation` looks more like this:

```
// After compilation, this function will be in-memory and usable as an
// Operation.

sixOrMore := func(in Value) Value {

    max := func(in Value) Value {

        a := tupEl(
            func(in Value) Value { return in }(in),
            func(_ Value) Value { return 0}(in),
        )

        b := tupEl(
            func(in Value) Value { return in }(in),
            func(_ Value) Value { return 1}(in),
        )

        if a > b {
            return a
        }

        return b
    }

    return max(
        func(in Value) Value { return in }(in),
        func(_ Value) Value { return 6}(in),
    )
}
```

This added layer of indirection for all leaf values is not great for
performance, and there's probably further refactoring which could be done to
make the result look more like the original ideal.

To make things a bit messier, even that representation isn't quite accurate to
the current result. The compiler doesn't properly de-duplicate work when
following name values. In other words, everytime `a` is referenced within `max`,
the `Operation` which the compiler produces will recompute `a` via `tupEl`.

So the _actual_ compiled `Operation` looks more like this:

```
// After compilation, this function will be in-memory and usable as an
// Operation.

sixOrMore := func(in Value) Value {

    return func(in Value) Value {

        if tupEl(func(in Value) Value { return in }(in), func(_ Value) Value { return 0}(in)) >
            tupEl(func(in Value) Value { return in }(in), func(_ Value) Value { return 1}(in)) {

            return tupEl(func(in Value) Value { return in }(in), func(_ Value) Value { return 0}(in))
        }

        return tupEl(func(in Value) Value { return in }(in), func(_ Value) Value { return 1}(in))
    }(
        func(in Value) Value { return in }(in),
        func(_ Value) Value { return 6}(in),
    )
}
```

Clearly, there's some optimization to be done still.

## Results

While it's still not perfect, the new implementation is far and away better than
the old. This can be seen just in the performance for the fibonacci program:

```
# Previous VM

$ time ./eval "$(cat examples/fib.gg)" 10
55

real    0m8.737s
user    0m9.871s
sys     0m0.309s
```

```
# New VM

$ time ./eval "$(cat examples/fib.gg)" 50
12586269025

real    0m0.003s
user    0m0.003s
sys     0m0.000s
```

They're not even comparable.
