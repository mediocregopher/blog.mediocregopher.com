---
title: >-
    Conditionals in Ginger
description: >-
    Some different options for how "if" statements could work.
series: ginger
tags: tech
---

In the [last ginger post][last] I covered a broad overview of how I envisioned
ginger would work as a language, but there were two areas where I felt there was
some uncertainty: conditionals and loops. In this post I will be focusing on
conditionals, and going over a couple of options for how they could work.

[last]: {% post_url 2021-01-09-ginger %}

## Preface

By "conditional" I'm referring to what programmers generally know as the "if"
statement; some mechanism by which code can do one thing or another based on
circumstances at runtime. Without some form of a conditional a programming
language is not Turing-complete and can't be used for anything interesting.

Given that it's uncommon to have a loop without some kind of a conditional
inside of it (usually to exit the loop), but it's quite common to have a
conditional with no loop in sight, it makes more sense to cover conditionals
before loops. Whatever decision is reached regarding conditionals will impact
how loops work, but not necessarily the other way around.

For the duration of this post I will be attempting to construct a simple
operation which takes two integers as arguments. If the first is less than
the second then the operation returns the addition of the two, otherwise the
operation returns the second subtracted from the first. In `go` this operation
would look like:

```go
func op(a, b int) int {
    if a < b {
        return a + b
    }
    return b - a
}
```

## Pattern 1: Branches As Inputs

The pattern I'll lay out here is simultaneously the first pattern which came to
me when trying to figure this problem out, the pattern which is most like
existing mainstream programming languages, and (in my opinion) the worst pattern
of the bunch. Here is what it looks like:

```
        in -lt-> } -if-> out
                 }
       in -add-> }
                 }
in -1-> }        }
in -0-> } -sub-> }

```

The idea here is that the operation `if` could take a 3-tuple whose elements
are, respectively: a boolean, and two other edges which won't be evaluated until
`if` is evaluated. If the boolean is true then `if` outputs the output of the
first edge (the second element in the tuple), and otherwise it will output the
value of the second edge.

This idea doesn't work for a couple reasons. The biggest is that, if there were
multiple levels of `if` statements, the structure of the graph grows out
_leftward_, whereas the flow of data is rightwards. For someone reading the code
to know what `if` will produce in either case they must first backtrack through
the graph, find the origin of that branch, then track that leftward once again
to the `if`.

The other reason this doesn't work is because it doesn't jive with any pattern
for loops I've come up with. This isn't evident from this particular example,
but consider what this would look like if either branch of the `if` needed to
loop back to a previous point in the codepath. If that's a difficult or
confusing task for you, you're not alone.

## Pattern 2: Pattern Matching

There's quite a few languages with pattern matching, and even one which I know
of (erlang) where pattern matching is the primary form of conditionals, and the
more common `if` statement is just some syntactic sugar on top of the pattern
matching.

I've considered pattern matching for ginger. It might look something like:

{% raw %}
```
       in -> } -switch-> } -> {{{A, B}, _}, ({A,B}-lt->out)} -0-> } -add-> out
in -1-> } -> }           } -1-> } -sub-> out
in -0-> }
```
{% endraw %}

The `switch` operation posits that a node can have multiple output edges. In a
graph this is fine, but it's worth noting. Graphs tend to be implemented such
that edges to and from a node are unordered, but in ginger it seems unlikely
that that will be the case.

The last output edge from the switch is the easiest to explain: it outputs the
input value to `switch` when no other branches are able to be taken. But the
input to `switch` is a bit complex in this example: It's a 2-tuple whose first
element is `in`, and whose second element is `in` but with reversed elements.
In the last output edge we immediately pipe into a `1` operation to retrieve
that second element and call `sub` on that, since that's the required behavior
of the example.

All other branches (in this switch there is only one, the first branch) output
to a value. The form of this value is a tuple (denoted by enclosed curly braces
here) of two values. The first value is the pattern itself, and the second is an
optional predicate. The pattern in this example will match a 2-tuple, ignoring
the second element in that tuple. The first element will itself be matched
against a 2-tuple, and assign each element to the variables `A` and `B`,
respectively. The second element in the tuple, the predicate, is a sub-graph
which returns a boolean, and can be used for further specificity which can't be
covered by the pattern matching (in this case, comparing the two values to each
other).

The output from any of `switch`'s branches is the same as its input value, the
only question is which branch is taken. This means that there's no backtracking
when reading a program using this pattern; no matter where you're looking you
will only have to keep reading rightward to come to an `out`.

There's a few drawbacks with this approach. The first is that it's not actually
very easy to read. While pattern matching can be a really nice feature in
languages that design around it, I've never seen it used in a LISP-style
language where the syntax denotes actual datastructures, and I feel that in such
a context it's a bit unwieldy. I could be wrong.

The second drawback is that pattern matching is not simple to implement, and I'm
not even sure what it would look like in a language where graphs are the primary
datastructure. In the above example we're only matching into a tuple, but how
would you format the pattern for a multi-node, multi-edge graph? Perhaps it's
possible. But given that any such system could be implemented as a macro on top
of normal `if` statements, rather than doing it the other way around, it seems
better to start with the simpler option.

(I haven't talked about it yet, but I'd like for ginger to be portable to
multiple backends (i.e. different processor architectures, vms, etc). If the
builtins of the language are complex, then doing this will be a difficult task,
whereas if I'm conscious of that goal during design I think it can be made to be
very simple. In that light I'd prefer to not require pattern matching to be a
builtin.)

The third drawback is that the input to the `switch` requires careful ordering,
especially in cases like this one where a different value is needed depending on
which branch is taken. I don't consider this to be a huge drawback, as
encourages good data design and is a common consideration in other functional
languages.

## Pattern 3: Branches As Outputs

Taking a cue from the pattern matching example, we can go back to `if` and take
advantage of multiple output edges being a possibility:

```
       in -> } -> } -if-> } -0-> } -add-> out
in -1-> } -> }    }       } -1-> } -sub-> out
in -0-> }         }
                  }
         in -lt-> }
```

It's not perfect, but I'd say this is the nicest of the three options so far.
`if` is an operation which takes a 2-tuple. The second element of the tuple is a
boolean, if the boolean is true then `if` passes the first element of its tuple
to the first branch, otherwise it passes it to the second. In this way `if`
becomes kind of like a fork in a train track: it accepts some payload (the first
element of its input tuple) and depending on conditions (the second element) it
directs the payload one way or the other.

This pattern retains the benefits of the pattern matching example, where one
never needs to backtrack in order to understand what is about to happen next,
while also being much more readable and simpler to implement. It also retains
one of the drawbacks of the pattern matching example, in that the inputs to `if`
must be carefully organized based on the needs of the output branches. As
before, I don't consider this to be a huge drawback.

There's other modifications which might be made to this `if` to make it even
cleaner, e.g. one could make it accept a 3-tuple, rather than a 2-tuple, in
order to supply differing values to be used depending on which branch is taken.
To me these sorts of small niceties are better left to be implemented as macros,
built on top of a simpler but less pleasant builtin.

## Fin

If you have other ideas around how conditionals might be done in a graph-based
language please [email me][email]; any and all contributions are welcome! One
day I'll get around to actually implementing some of ginger, but today is not
that day.

[email]: mailto:mediocregopher@gmail.com
