---
title: >-
    Program Structure and Composability
description: >-
    Discussing the nature of program structure, the problems presented by
    complex structures, and a pattern which helps in solving those problems.
---

## Part 0: Intro

This post is focused on a concept I call "program structure", which I will try
to shed some light on before moving on to discussing complex program structures,
discussing why complex structures can be problematic to deal with, and finally
discussing a pattern for dealing with those problems.

My background is as a backend engineer working on large projects that have had
many moving parts; most had multiple services interacting, used many different
databases in various contexts, and faced large amounts of load from millions of
users. Most of this post will be framed from my perspective, and present
problems in the way I have experienced them. I believe, however, that the
concepts and problems I discuss here are applicable to many other domains, and I
hope those with a foot in both backend systems and a second domain can help to
translate the ideas between the two.

## Part 1: Program Structure

For a long time I thought about program structure in terms of the hierarchy
present in the filesystem. In my mind, a program's structure looked like this:

```
// The directory structure of a project called gobdns.
src/
    config/
    dns/
    http/
    ips/
    persist/
    repl/
    snapshot/
    main.go
```

What I grew to learn was that this consolidation of "program structure" with
"directory structure" is ultimately unhelpful. While I won't deny that every
program has a directory structure (and if not, it ought to), this does not mean
that the way the program looks in a filesystem in anyway corresponds to how it
looks in our mind's eye.

The most notable way to show this is to consider a library package. Here is the
structure of a simple web-app which uses redis (my favorite database) as a
backend:

```
src/
    redis/
    http/
    main.go
```

(Note that I use go as my example language throughout this post, but none of the
ideas I'll referring to are go specific.)

If I were to ask you, based on that directory strucure, what the program does,
in the most abstract terms, you might say something like: "The program
establishes an http server which listens for requests, as well as a connection
to the redis server. The program then interacts with redis in different ways,
based on the http requests which are received on the server."

And that would be a good guess. But consider another case: "The program
establishes an http server which listens for requests, as well as connections to
_two different_ redis servers. The program then interacts with one redis server
or the other in different ways, based on the http requests which are received
from the server.

The directory structure could apply to either description; `redis` is just a
library which allows for interacting with a redis server, but it doesn't specify
_which_ server, or _how many_. And those are extremely important factors which
are definitely reflected in our concept of the program's structure, and yet not
in the directory structure. Even worse, thinking of structure in terms of
directories might (and, I claim, often does) cause someone to assume that
program only _could_ interact with one redis server, which is obviously untrue.

### Global State and Microservices

The directory-centric approach to structure often leads to the use of global
singletons to manage access to external resources like RPC servers and
databases. In the above example the `redis` library might contain code which
looks something like:

```go
// For the non-gophers, redisConnection is variable type which has been made up
// for this example.
var globalConn redisConnection

func Get() redisConnection {
    if globalConn == nil {
        globalConn = makeConnection()
    }
    return globalConn
}
```

Ignoring that the above code is not thread-safe, the above pattern has some
serious drawbacks. For starters, it does not play nicely with a microservices
oriented system, or any other system with good separation of concerns between
its components.

I have been a part of building several large products with teams of various
sizes. In each case we had a common library which was shared amongst all
components of the system, and contained functionality which was desired to be
kept the same across those components. For example, configuration was generally
done through that library, so all components could be configured in the same
way. Similarly, an RPC framework is usually included in the common library, so
all components can communicate in a shared language. The common library also
generally contains domain specific types, for example a `User` type which all
components will need to be able to understand.

Most common libraries also have parts dedicated to databases, such as the
`redis` library example we've been using. In a medium-to-large sized system,
with many components, there are likely to be multiple running instances of any
database: multiple SQLs, different caches for each, different queues set up for
different asynchronous tasks, etc... And this is good! The ideal
compartmentalized system has components interact with each other directly, not
via their databases, and so each component ought to, to the extent possible,
keep its own databases to itself, with other components not touching them.

The singleton pattern breaks this separation, by forcing the configuration of
_all_ databases through the common library. If one component in the system adds
a database instance, all other components have access to it. While this doesn't
necessarily mean the components will _use_ it, that will only be accomplished
through sheer discipline, which will inevitably break down once management
decides it's crunch time.

To be clear, I'm not suggesting that singletons make proper compartmentalization
impossible, they simply add friction to it. In other words, compartmentalization
is not the default mode of singletons.

Another problem with singletons, as mentioned before, is that they don't handle
multiple instances of the same thing very well. In order to support having
multiple redis instances in the system, the above code would need to be modified
to give every instance a name, and track the mapping of between that name, its
singleton, and its configuration. For large projects the number of different
instances can be enormous, and often the list which exists in code does not stay
fully up-to-date.

This might all sound petty, but I think it has a large impact. Ultimately, when
a component is using a singleton which is housed in a common library, that
component is borrowing the instance, rather than owning it. Put another way, the
component's structure is partially held by the common library, and since all
components are going to use the common library, all of their structures are
incorporated together. The separation between components is less solidified, and
systems become weaker.

What I'm going to propose is an alternative way to think about program structure
which still allows for all the useful aspects of a common library, without
compromising on component separation, and therefore giving large teams more
freedom to act independently of each other.
