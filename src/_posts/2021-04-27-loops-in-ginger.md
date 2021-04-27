---
title: >-
    Loops in Ginger
description: >-
    Bringing it back around.
series: ginger
tags: tech
---

In previous posts in this series I went over the general idea of the ginger
programming language, and some of its properties. To recap:

* Ginger is a programming language whose syntax defines a directed graph, in the
  same way that a LISP language's syntax defines nested lists.

* Graph edges indicate an operation, while nodes indicate a value.

* The special values `in` and `out` are used when interpreting a graph as a
  function.

* A special node type, the tuple, is defined as being a node whose value is an
  ordered set of input edges.

* Another special node type, the fork, is the complement to the tuple. A fork is
  defined as being a node whose value is an ordered set of output edges.

* The special `if` operation accepts a 2-tuple, the first value being some state
  value and the second being a tuple. The `if` operation expects to be directed
  towards a 2-fork. If the boolean is true then the top output edge of the fork
  is taken, otherwise the bottom is taken. The state value is what's passed to
  the taken edge.

There were some other detail rules but I don't remember them off the top of my
head.

## Loops

Today I'd like to go over my ideas for how loops would work in ginger. With
loops established ginger would officially be a Turing complete language and,
given time and energy, real work could actually begin on it.

As with conditionals I'll start by establishing a base example. Let's say we'd
like to define an operation which prints out numbers from 0 up to `n`, where `n`
is given as an argument. In go this would look like:

```go
func printRange(n int) int {
    for i := 0; i < n; i++ {
        fmt.Println(i)
    }
}
```

With that established, let's start looking at different patterns.

## Goto

In the olden days the primary looping construct was `goto`, which essentially
teleports the program counter (aka instruction pointer) to another place in the
execution stack. Pretty much any other looping construct can be derived from
`goto` and some kind of conditional, so it's a good starting place when
considering loops in ginger.

```
(in -println-> } -incr-> out) -> println-incr

0  -> }    -> } -if-> { -> out
in -> } -eq-> }       { -> } -upd-> } -+
      ^               0 -> }           |
      |    println-incr -> }           |
      |                                |
      +--------------------------------+
```

(Note: the `upd` operation is used here for convenience. It takes in three
arguments: A tuple, an index, and an operation. It applies the operation to the
tuple element at the given index, and returns a new tuple with that index set to
the value returned.)

Here `goto` is performed using a literal arrow going from the right to left.
it's ugly and hard to write, and would only be moreso the more possible gotos an
operation has.

It also complicates our graphs in a significant way: up till now ginger graphs
have have always been directed _acyclic_ graphs (DAGs), but by introducing this
construct we allow that graphs might be cyclic. It's not immediately clear to me
what the consequences of this will be, but I'm sure they will be great. If
nothign else it will make the compiler much more complex, as each value can no
longer be defined in terms of its input edge, as that edge might resolve back to
the value itself.

While conceptually sound, I think this strategy fails the practability test. We
can do better.

## While

The `while` construct is the basic looping primitive of iterative languages
(some call it `for`, but they're just lying to themselves).

Try as I might, I can't come up with a way to make `while` work with ginger.
`while` ultimately relies on scoped variables being updated in place to
function, while ginger is based on the concept of pipelining a set of values
through a series of operations. From the point of view of the programmer these
operations are essentially immutable, so the requirement of a variable which can
be updated in place cannot be met.

## Recur

This pattern is based on how many functional languages, for example erlang,
handle looping. Rather than introducing new primitives around looping, these
language instead ensure that tail calls are properly optimized and uses those
instead. So loops are implemented as recursive function calls.

For ginger to do this it would make sense to introduce a new special value,
`recur`, which could be used alongside `in` and `out` within operations. When
the execution path hits a `recur` then it gets teleported back to the `in`
value, with the input to `recur` now being the output from `in`. Usage of it
would look like:

```
(

    (in -println-> } -incr-> out) -> println-incr

    in    -> } -if-> { -> out
    in -eq-> }       { -> } -upd-> } -> recur
                     0 -> }
          println-incr -> }

) -> inner-op

0  -> } -inner-op-> out
in -> }
```

This looks pretty similar to the `goto` example overall, but with the major
difference that the looping body had to be wrapped into an inner operation. The
reason for this is that the outer operation only takes in one argument, `n`, but
the loop actually needs two pieces of state to function: `n` and the current
value. So the inner operation loops over these two pieces of state, and the
outer operation supplies `n` and an initial iteration value (`0`) to that inner
operation.

This seems cumbersome on the surface, but what other languages do (such as
erlang, which is the one I'm most familiar with) is to provide built-in macros
on top of this primitive which make it more pleasant to use. These include
function polymorphism and a more familiar `for` construct. With a decent macro
capability ginger could do the same.

The benefits here are that the graphs remain acyclic, and the syntax has not
been made more cumbersome. It follows conventions established by other
languages, and ensures the language will be capable of tail-recursion.

## Map/Reduce

Another functional strategy which is useful is that of the map/reduce power
couple. The `map` operation takes a sequence of values and an operation, and
returns a sequence of the same length where the operation has been applied to
each value in the original sequence individually. The `reduce` operation is more
complicated (and not necessary for out example), but it's essentially a
mechanism to turn a sequence of values into a single value.

For our example we only need `map`, plus one more helper operation: `range`.
`range` takes a number `n` and returns a sequence of numbers starting at `0` and
ending at `n-1`. Our print example now looks like:

```
in -range-> } -map-> out
 println -> }
```

Very simple! Map/reduce is a well established pattern and is probably the
best way to construct functional programs. However, the question remains whether
these are the best _primitives_ for looping, and I don't believe they are. Both
`map` and `reduce` can be derived from conditional and looping primitives like
`if` and `recur`, and they can't do some things that those primitives can. While


I expect one of the first things which will be done in ginger is to define `map`
and `reduce` in terms of `if` and a looping primitive, and use them generously
throughout the code, I think the fact that they can be defined in terms of
lower-level primitives indicates that they aren't the right looping primitives
for ginger.

## Conclusion

Unlike with the conditionals posts, where I started out not really knowing what
I wanted to do with conditionals, I more or less knew where this post was going
from the beginning. `recur` is, in my mind, the best primitive for looping in
ginger. It provides the flexibility to be extended to any use-case, while not
complicating the structure of the language. While possibly cumbersome to
implement directly, `recur` can be used as a primitive to construct more
convenient looping operations like `map` and `reduce`.

As a final treat (lucky you!), here's `map` defined using `if` and `recur`:

```
(
    in -0-> mapped-seq
    in -1-> orig-seq
    in -2-> op

    mapped-seq -len-> i

              mapped-seq -> } -if-> { -> out
    orig-seq -len-> } -eq-> }       { -> } -append-> } -> recur
               i -> }                    }           }
                                         }           }
                   orig-seq -i-> } -op-> }           }
                                                     }
                                         orig-seq -> }
                                               op -> }
) -> inner-map

  () -> } -inner-map-> out
in -0-> }
in -1-> }
```

The next step for ginger is going to be writing an actual implementation of the
graph structure in some other language (let's be honest, it'll be in go). After
that we'll need a syntax definition which can be used to encode/decode that
structure, and from there we can start actually implementing the language!
