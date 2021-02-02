---
title: >-
    Component-Oriented Programming
description: >-
    A concise description of.
tags: tech
---

[A previous post in this
blog](/2019/08/02/program-structure-and-composability.html) focused on a
framework developed to make designing component-based programs easier. In
retrospect, the proposed pattern/framework was over-engineered. This post
attempts to present the same ideas in a more distilled form, as a simple
programming pattern and without the unnecessary framework.

## Components

Many languages, libraries, and patterns make use of a concept called a
"component," but in each case the meaning of "component" might be slightly
different. Therefore, to begin talking about components, it is necessary to first
describe what is meant by "component" in this post.

For the purposes of this post, the properties of components include the
following.

&nbsp;1... **Abstract**: A component is an interface consisting of one or more
methods. 

&nbsp;&nbsp;&nbsp;1a... A function might be considered a single-method component
_if_ the language supports first-class functions.

&nbsp;&nbsp;&nbsp;1b... A component, being an interface, may have one or more
implementations. Generally, there will be a primary implementation, which is
used during a program's runtime, and secondary "mock" implementations, which are
only used when testing other components.

&nbsp;2... **Instantiatable**: An instance of a component, given some set of
parameters, can be instantiated as a standalone entity. More than one of the
same component can be instantiated, as needed.

&nbsp;3... **Composable**: A component may be used as a parameter of another
component's instantiation. This would make it a child component of the one being
instantiated (the parent).

&nbsp;4... **Pure**: A component may not use mutable global variables (i.e.,
singletons) or impure global functions (e.g., system calls). It may only use
constants and variables/components given to it during instantiation.

&nbsp;5... **Ephemeral**: A component may have a specific method used to clean
up all resources that it's holding (e.g., network connections, file handles,
language-specific lightweight threads, etc.).

&nbsp;&nbsp;&nbsp;5a... This cleanup method should _not_ clean up any child
components given as instantiation parameters.

&nbsp;&nbsp;&nbsp;5b... This cleanup method should not return until the
component's cleanup is complete.

&nbsp;&nbsp;&nbsp;5c... A component should not be cleaned up until all its
parent components are cleaned up.

Components are composed together to create component-oriented programs. This is
done by passing components as parameters to other components during
instantiation. The `main` procedure of the program is responsible for
instantiating and composing the components of the program.

## Example

It's easier to show than to tell. This section posits a simple program and then
describes how it would be implemented in a component-oriented way. The program
chooses a random number and exposes an HTTP interface that allows users to try
and guess that number. The following are requirements of the program:

* A guess consists of a name that identifies the user performing the guess and
  the number that is being guessed;

* A score is kept for each user who has performed a guess;

* Upon an incorrect guess, the user should be informed of whether they guessed
  too high or too low, and 1 point should be deducted from their score;

* Upon a correct guess, the program should pick a new random number against
  which to check subsequent guesses, and 1000 points should be added to the
  user's score;

* The HTTP interface should have two endpoints: one for users to submit guesses,
  and another that lists out user scores from highest to lowest;

* Scores should be saved to disk so they survive program restarts.

It seems clear that there will be two major areas of functionality for our
program: score-keeping and user interaction via HTTP. Each of these can be
encapsulated into components called `scoreboard` and `httpHandlers`,
respectively.

`scoreboard` will need to interact with a filesystem component to save/restore
scores (because it can't use system calls directly; see property 4). It would be
wasteful for `scoreboard` to save the scores to disk on every score update, so
instead it will do so every 5 seconds. A time component will be required to
support this.

`httpHandlers` will be choosing the random number which is being guessed, and
will therefore need a component that produces random numbers. `httpHandlers`
will also be recording score changes to `scoreboard`, so it will need access to
`scoreboard`.

The example implementation will be written in go, which makes differentiating
HTTP handler functionality from the actual HTTP server quite easy; thus, there
will be an `httpServer` component that uses `httpHandlers`.

Finally, a `logger` component will be used in various places to log useful
information during runtime.

[The example implementation can be found
here.](/assets/component-oriented-design/v1/main.html) While most of it can be
skimmed, it is recommended to at least read through the `main` function to see
how components are composed together. Note that `main` is where all components
are instantiated, and that all components' take in their child components as
part of their instantiation.

## DAG

One way to look at a component-oriented program is as a directed acyclic graph
(DAG), where each node in the graph represents a component, and each edge
indicates that one component depends upon another component for instantiation.
For the previous program, it's quite easy to construct such a DAG just by
looking at `main`, as in the following:

```
net.Listener     rand.Rand        os.File
     ^               ^               ^
     |               |               |
 httpServer --> httpHandlers --> scoreboard --> time.Ticker
     |               |               |
     +---------------+---------------+--> log.Logger
```

Note that all the leaves of the DAG (i.e., nodes with no children) describe the
points where the program meets the operating system via system calls. The leaves
are, in essence, the program's interface with the outside world.

While it's not necessary to actually draw out the DAG for every program one
writes, it can be helpful to at least think about the program's structure in
these terms.

## Benefits

Looking at the previous example implementation, one would be forgiven for having
the immediate reaction of "This seems like a lot of extra work for little gain.
Why can't I just make the system calls where I need to, and not bother with
wrapping them in interfaces and all these other rules?"

The following sections will answer that concern by showing the benefits gained
by following a component-oriented pattern.

### Testing

Testing is important, that much is being assumed.

A distinction to be made with testing is between unit and non-unit tests. Unit
tests are those for which there are no requirements for the environment outside
the test, such as the existence of global variables, running databases,
filesystems, or network services. Unit tests do not interact with the world
outside the testing procedure, but instead use mocks in place of the
functionality that would be expected by that world.

Unit tests are important because they are faster to run and more consistent than
non-unit tests. Unit tests also force the programmer to consider different
possible states of a component's dependencies during the mocking process.

Unit tests are often not employed by programmers, because they are difficult to
implement for code that does not expose any way to swap out dependencies for
mocks of those dependencies. The primary culprit of this difficulty is the
direct usage of singletons and impure global functions. For component-oriented
programs, all components inherently allow for the swapping out of any
dependencies via their instantiation parameters, so there's no extra effort
needed to support unit tests.

[Tests for the example implementation can be found
here.](/assets/component-oriented-design/v1/main_test.html) Note that all
dependencies of each component being tested are mocked/stubbed next to them.

### Configuration

Practically all programs require some level of runtime configuration. This may
take the form of command-line arguments, environment variables, configuration
files, etc.

For a component-oriented program, all components are instantiated in the same
place, `main`, so it's very easy to expose any arbitrary parameter to the user
via configuration. For any component that is affected by a configurable
parameter, that component merely needs to take an instantiation parameter for
that configurable parameter; `main` can connect the two together. This accounts
for the unit testing of a component with different configurations, while still
allowing for the configuration of any arbitrary internal functionality.

For more complex configuration systems, it is also possible to implement a
`configuration` component that wraps whatever configuration-related
functionality is needed, which other components use as a sub-component. The
effect is the same.

To demonstrate how configuration works in a component-oriented program, the
example program's requirements will be augmented to include the following:

* The point change values for both correct and incorrect guesses (currently
  hardcoded at 1000 and 1, respectively) should be configurable on the
  command-line;

* The save file's path, HTTP listen address, and save interval should all be
  configurable on the command-line.

[The new implementation, with newly configurable parameters, can be found
here.](/assets/component-oriented-design/v2/main.html) Most of the program has
remained the same, and all unit tests from before remain valid. The primary
difference is that `scoreboard` takes in two new parameters for the point change
values, and configuration is set up inside `main` using the `flags` package.

### Setup/Runtime/Cleanup

A program can be split into three stages: setup, runtime, and cleanup. Setup is
the stage during which the internal state is assembled to make runtime possible.
Runtime is the stage during which a program's actual function is being
performed. Cleanup is the stage during which the runtime stops and internal
state is disassembled.

A graceful (i.e., reliably correct) setup is quite natural to accomplish for
most. On the other hand, a graceful cleanup is, unfortunately, not a programmer's
first concern (if it is a concern at all).

When building reliable and correct programs, a graceful cleanup is as important
as a graceful setup and runtime. A program is still running while it is being
cleaned up, and it's possibly still acting on the outside world. Shouldn't
it behave correctly during that time?

Achieving a graceful setup and cleanup with components is quite simple.

During setup, a single-threaded procedure (`main`) first constructs the leaf
components, then the components that take those leaves as parameters, then the
components that take _those_ as parameters, and so on, until the component DAG
is fully constructed.

At this point, the program's runtime has begun.

Once the runtime is over, signified by a process signal or some other mechanism,
it's only necessary to call each component's cleanup method (if any; see
property 5) in the reverse of the order in which the components were
instantiated.  This order is inherently deterministic, as the components were
instantiated by a single-threaded procedure.

Inherent to this pattern is the fact that each component will certainly be
cleaned up before any of its child components, as its child components must have
been instantiated first, and a component will not clean up child components
given as parameters (properties 5a and 5c). Therefore, the pattern avoids
use-after-cleanup situations.

To demonstrate a graceful cleanup in a component-oriented program, the example
program's requirements will be augmented to include the following:

* The program will terminate itself upon an interrupt signal;

* During termination (cleanup), the program will save the latest set of scores
  to disk one final time.

[The new implementation that accounts for these new requirements can be found
here.](/assets/component-oriented-design/v3/main.html) For this example, go's
`defer` feature could have been used instead, which would have been even
cleaner, but was omitted for the sake of those using other languages.


## Conclusion

The component pattern helps make programs more reliable with only a small amount
of extra effort incurred. In fact, most of the pattern has to do with
establishing sensible abstractions around global functionality and remembering
certain idioms for how those abstractions should be composed together, something
most of us already do to some extent anyway.

While beneficial in many ways, component-oriented programming is merely a tool
that can be applied in many cases. It is certain that there are cases where it
is not the right tool for the job, so apply it deliberately and intelligently.

## Criticisms/Questions

In lieu of a FAQ, I will attempt to premeditate questions and criticisms of the
component-oriented programming pattern laid out in this post.

**This seems like a lot of extra work.**

Building reliable programs is a lot of work, just as building a
reliable _anything_ is a lot of work. Many of us work in an industry that likes
to balance reliability (sometimes referred to by the more specious "quality")
with malleability and deliverability, which naturally leads to skepticism of any
suggestions requiring more time spent on reliability. This is not necessarily a
bad thing, it's just how the industry functions.

All that said, a pattern need not be followed perfectly to be worthwhile, and
the amount of extra work incurred by it can be decided based on practical
considerations. I merely maintain that code which is (mostly) component-oriented
is easier to maintain in the long run, even if it might be harder to get off the
ground initially.

**My language makes this difficult.**

I don't know of any language which makes this pattern particularly easier than
others, so, unfortunately, we're all in the same boat to some extent (though I
recognize that some languages, or their ecosystems, make it more difficult than
others). It seems to me that this pattern shouldn't be unbearably difficult for
anyone to implement in any language either, however, as the only language
feature required is abstract typing.

It would be nice to one day see a language that explicitly supports this
pattern by baking the component properties in as compiler-checked rules.

**My `main` is too big**

There's no law saying all component construction needs to happen in `main`,
that's just the most sensible place for it. If there are large sections of your
program that are independent of each other, then they could each have their own
construction functions that `main` then calls.

Other questions that are worth asking include: Can my program be split up
into multiple programs? Can the responsibilities of any of my components be
refactored to reduce the overall complexity of the component DAG? Can the
instantiation of any components be moved within their parent's
instantiation function?

(This last suggestion may seem to be disallowed, but is fine as long as the
parent's instantiation function remains pure.)

**Won't this will result in over-abstraction?**

Abstraction is a necessary tool in a programmer's toolkit, there is simply no
way around it. The only questions are "how much?" and "where?"

The use of this pattern does not affect how those questions are answered, in my
opinion, but instead aims to more clearly delineate the relationships and
interactions between the different abstracted types once they've been
established using other methods. Over-abstraction is possible and avoidable
regardless of which language, pattern, or framework is being used.

**Does CoP conflict with object-oriented or functional programming?**

I don't think so. OoP languages will have abstract types as part of their core
feature-set; most difficulties are going to be with deliberately _not_ using
other features of an OoP language, and with imported libraries in the language
perhaps making life inconvenient by not following CoP (specifically regarding
cleanup and the use of singletons).

For functional programming, it may well be that, depending on the language, CoP
is technically being used, as functional languages are already generally
antagonistic toward globals and impure functions, which is most of the battle.
If anything, the transition from functional to component-oriented programming
will generally be an organizational task.
