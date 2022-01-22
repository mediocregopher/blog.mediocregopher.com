---
title: >-
    Ginger Names
description: >-
    Thoughts about a fundamental data type.
tags: tech
series: ginger
---

The ginger language has, so far, 2 core types implemented: numbers and names.
Obviously there will be more coming later, but at this stage of development
these are all that's really needed. Numbers are pretty self explanatory, but
it's worth talking about names a bit.

As they are currently defined, a name's only property is that it can either be
equal or not equal to another name. Syntactically they are encoded as being any
alphanumeric token starting with an alphabetic character. We might _think_ of
them as being strings, but names lack nearly all capabilities that strings have:
they cannot be iterated over, they cannot be concatenated, they cannot be split.
Names can only be compared for equality.

## Utility

The use-case for names is self-explanatory: they are words which identify
something from amongst a group.

Consider your own name. It _might_ have an ostensible meaning. Mine, Brian,
means "high" (as in... like a hill, which is the possible root word). But when
people yell "Brian" across the room I'm in, they don't mean a hill. They mean
me, because that word is used to identify me from amongst others. The etymology
is essentially background information which doesn't matter.

We use names all the time in programming, though we don't always call them that.
Variable names, package names, type names, function names, struct field names.
There's also keys which get used in hash maps, which are essentially names, as
well as enumerations. By defining name as a core type we can cover a lot of
ground.

## Precedence

This is not the first time a name has been used as a core type. Ruby has
symbols, which look like `:this`. Clojure has keywords, which also look like
`:this`, and it has symbols, which look like `this`. Erlang has atoms, which
don't have a prefix and so look like `this`. I can't imagine these are the only
examples. They are all called different things, but they're all essentially the
same thing: a runtime value which can only be compared for equality.

I can't speak much about ruby, but I _can_ speak about clojure and erlang.

Clojure is a LISP language, meaning the language itself is described using the
data types and structures built into the language. Ginger is also a LISP, though
it uses graphs instead of lists.

Clojure keywords are generally used as keys to hash maps, sentinel values, and
enumerations. Besides keywords, clojure also makes use of symbols, which are
used for variable and library names. There seems to be some kind of split
ability on symbols, as they are expected to be separated on their periods when
importing, as in `clojure.java.io`. There's also a quoting mechanism in clojure,
where prefixing a symbol, or other value, with a single quote, like `'this`,
prevents it from being evaluated as a variable or function call.

It's also possible to have something get quoted multiple layers deep, like
`'''this`. This can get confusing.

Erlang is not a LISP language, but it does have atoms. These values are used in
the same way that clojure keywords are used. There is no need for a
corresponding symbol type like clojure has, since erlang is not a LISP and has
no real macro system. Atoms are sort of used like symbols, in that functions and
packages are identified by an atom, and so one can "call" an atom, like
`this()`, in order to evaluate it.

## Just Names

I don't really see the need for clojure's separation between keywords and
symbols. Symbols still need to be quoted in order to prevent evaluation either
way, so you end up with three different entities to juggle (keywords, symbols,
and symbols which won't be evaluated). Erlang's solution is simpler, atoms are
just atoms, and since evaluation is explicit there's no need for quoting. Ginger
names are like erlang atoms in that they are the only tool at hand.

The approaches of erlang vs clojure could be reframed as explicit vs implicit
evaluation of operations calls.

In ginger evaluation is currently done implicitly, but in only two cases:

* A value on an edge is evaluated to the first value which is a graph (which
  then gets interpreted as an operation).

* A leaf vertex with a name value is evaluated to the first value which is not a
  name.

In all other cases, the value is left as-is. A graph does not need to be quoted,
since the need to evaluate a graph as an operation is already based on its
placement as an edge or not. So the only case left where quoting is needed (if
implicit evaluation continues to be used) is a name on a leaf vertex, as in the
example before.

As an example to explore explicit vs implicit quoting in ginger, if we want to
programatically call the `AddValueIn` method on a graph, which terminates an
open edge into a value, and that value is a name, it might look like this with
implicit evaluation (the clojure-like example):

```
out = addValueIn < (g (quote < someName;) someValue; );

* or, to borrow the clojure syntax, where single quote is a shortcut:

out = addValueIn < (g; 'someName; someValue; );
```

In an explicit evaluation language, which ginger so far has not been and so this
will look weird, we might end up with something like this:

```
out = addValueIn < (eval < g; someName; eval < someValue; );

* with $ as sugar for the `eval`, like ' is a shortcut for `quote` in clojure:`

out = addValueIn < ($g; someName; $someValue; );
```

I don't _like_ either pattern, and since it's such a specific case I feel like
something less obtrusive could come up. So no decisions here yet.

## Uniqueness

There's another idea I haven't really gotten to the bottom of yet. The idea is
that a name, _maybe_, shouldn't be considered equal to the same name unless they
belong to the same graph.

For example:

```
otherFoo = { out = 'foo } < ();

out = equal < ('foo;  otherFoo; );
```

This would output false. `otherFoo`'s value is the name `foo`, and the value
it's being compared to is also a name `foo`, but they are from different graphs
and so are not equal. In essence, names are automatically namespaces.

This idea only really makes sense in the context of packages, where a user
(a developer) wants to import functionality from somewhere else and use it
in their program. The code package which is imported will likely use name
values internally to implement its functionality, but it shouldn't need to worry
about naming conflicts with values passed in by the user. While it's possible to
avoid conflicts if a package is designed conscientiously, it's also easy to mess
up if one isn't careful. This becomes especially true when combining
functionality of packages with overlapping functionality, where the data
returned from one might looks _similar_ to that used by the other, but it's not
necessarily true.

On the other hand, this could create some real headaches for the developer, as
they chase down errors which are caused because one `foo` isn't actually the
same as another `foo`.

What it really comes down to is the mechanism which packages use to function as
packages. Forced namespaces will require packages to export all names which they
expect the user to need to work with the package. So the ergonomics of that
exporting, both on the user's and package's side, are really important in order
to make this bearable.

So it's hard to make any progress on determining if this idea is gonna work
until the details of packaging are worked out. But for this idea to work the
packaging is going to need to be designed with it in mind. It's a bit of a
puzzle, and one that I'm going to marinate on longer, in addition to the quoting
of names.

And that's names, their current behavior and possible future behavior. Keep an
eye out for more ginger posts in.... many months, because I'm going to go work
on other things for a while (I say, with a post from a month ago having ended
with the same sentiment).
