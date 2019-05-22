---
title: >-
    Program Structure and Composability
description: >-
    Discussing the nature of program structure, the problems presented by
    complex structures, and a pattern which helps in solving those problems.
---

TODO:
* Double check if I'm using "I" or "We" everywhere (probably should use "I")

## Part 0: Introduction

This post is focused on a concept I call "program structure", which I will try
to shed some light on before discussing complex program structures, then
discussing why complex structures can be problematic to deal with, and finally
discussing a pattern for dealing with those problems.

My background is as a backend engineer working on large projects that have had
many moving parts; most had multiple services interacting with each other, using
many different databases in various contexts, and facing large amounts of load
from millions of users. Most of this post will be framed from my perspective,
and will present problems in the way I have experienced them. I believe,
however, that the concepts and problems I discuss here are applicable to many
other domains, and I hope those with a foot in both backend systems and a second
domain can help to translate the ideas between the two.

Also note that I will be using Go as my example language, but none of the
concepts discussed here are specific to Go. To that end, I've decided to favor
readable code over "correct" code, and so have elided things that most gophers
hold near-and-dear, such as error checking and comments on all public types, in
order to make the code as accessible as possible to non-gophers as well. As with
before, I trust someone with a foot in Go and another language can translate
help me translate between the two.

## Part 1: Program Structure

In this section I will discuss the difference between directory and program
structure, show how global state is antithetical to compartmentalization (and
therefore good program structure), and finally discuss a more effective way to
think about program structure.

### Directory Structure

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

What I grew to learn was that this conflation of "program structure" with
"directory structure" is ultimately unhelpful. While I won't deny that every
program has a directory structure (and if not, it ought to), this does not mean
that the way the program looks in a filesystem in any way corresponds to how it
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

If I were to ask you, based on that directory strucure, what the program does,
in the most abstract terms, you might say something like: "The program
establishes an http server which listens for requests, as well as a connection
to the redis server. The program then interacts with redis in different ways,
based on the http requests which are received on the server."

And that would be a good guess. Here's a diagram which depicts the program
structure, wherein the root node, `main.go`, takes in requests from `http` and
processes them using `redis`.

TODO diagram

This is certainly a viable guess for how a program with that directory structure
operates, but consider another: "A component of the program called `server`
establishes an http server which listens for requests, as well as a connection
to a redis server. `server` then interacts with that redis connection in
different ways, based on the http requests which are received on the http
server.  Additionally, `server` tracks statistics about these interactions and
makes them available to other components. The root component of the program
establishes a connection to a second redis server, and stores those statistics
in that redis server."

TODO diagram

The directory structure could apply to either description; `redis` is just a
library which allows for interacting with a redis server, but it doesn't specify
_which_ server, or _how many_. And those are extremely important factors which
are definitely reflected in our concept of the program's structure, and yet not
in the directory structure. **What the directory structure reflects are the
different _kinds_ of components available to use, but it does not reflect how a
program will use those components.**

### Global State vs. Compartmentalization

The directory-centric approach to structure often leads to the use of global
singletons to manage access to external resources like RPC servers and
databases. In the above example the `redis` library might contain code which
looks something like:

```go
// A mapping of connection names to redis connections.
var globalConns = map[string]*RedisConn{}

func Get(name string) *RedisConn {
    if globalConns[name] == nil {
        globalConns[name] = makeConnection(name)
    }
    return globalConns[name]
}
```

Even though this pattern would work, it breaks with our conception of the
program structure in the more complex case shown above. Rather than having the
`server` component own the redis server it uses, the root component would be the
owner of it, and `server` would be borrowing it. Compartmentalization has been
broken, and can only be held together through sheer human discipline.

This is the problem with all global state. It's shareable amongst all components
of a program, and so is owned by none of them. One must look at an entire
codebase to understand how a globally held component is used, which might not
even be possible for a large codebase. And so the maintainers of these shared
components rely entirely on the discipline of their fellow coders when making
changes, usually discovering where that discipline broke down once the changes
have been pushed live.

Global state also makes it easier for disparate services/components to share
datastores for completely unrelated tasks. In the above example, rather than
creating a new redis instance for the root component's statistics storage, the
coder might have instead said "well, there's already a redis instance available,
I'll just use that." And so compartmentalization would have been broken further.
Perhaps the two instances _could_ be coalesced into the same one, for the sake
of resource efficiency, but that decision would be better made at runtime via
the configuration of the program, rather than being hardcoded into the code.

From the perspective of team management, global state-based patterns do nothing
except slow teams down. The person/team responsible for maintaining the central
library which holds all the shared resources (`redis`, in the above example)
becomes the bottleneck for creating new instances for new components, which will
further lead to re-using existing instances rather than create new ones, further
breaking compartmentalization. The person/team responsible for the central
library often finds themselves as the maintainers of the shared resource as
well, rather than the team actually using it.

### Component Structure

So what does proper program structure look like? In my mind the structure of a
program is a hierarchy of components, or, in other words, a tree. The leaf nodes
of the tree are almost _always_ IO related components, e.g. database
connections, RPC server frameworks or clients, message queue consumers, etc...
The non-leaf nodes will _generally_ be components which bring together the
functionalities of their children in some useful way, though they may also have
some IO functionality of their own.

Let's look at an even more complex structure, still only using the `redis` and
`http` component types:

TODO diagram:
```
    root
        rest-api
            redis
            http
        redis // for stats keeping
        debug
            http
```

This component structure contains the addition of the `debug` component. Clearly
the `http` and `redis` components are reusable in different contexts, but for
this example the `debug` endpoint is as well. It creates a separate http server
which can be queried to perform runtime debugging of the program, and can be
tacked onto virtually any program. The `rest-api` component is specific to this
program and therefore not reusable. Let's dive into it a bit to see how it might
be implemented:

```go
// RestAPI is very much not thread-safe, hopefully it doesn't have to handle
// more than one request at once.
type RestAPI struct {
    redisConn *redis.RedisConn
    httpSrv   *http.Server

    // Statistics exported for other components to see
    RequestCount int
    FooRequestCount int
    BarRequestCount int
}

func NewRestAPI() *RestAPI {
    r := new(RestAPI)
    r.redisConn := redis.NewConn("127.0.0.1:6379")

    // mux will route requests to different handlers based on their URL path.
    mux := http.NewServeMux()
    mux.Handle("/foo", http.HandlerFunc(r.fooHandler))
    mux.Handle("/bar", http.HandlerFunc(r.barHandler))
    r.httpSrv := http.NewServer(mux)

    // Listen for requests and serve them in the background.
    go r.httpSrv.Listen(":8000")

    return r
}

func (r *RestAPI) fooHandler(rw http.ResponseWriter, r *http.Request) {
    r.redisConn.Command("INCR", "fooKey")
    r.RequestCount++
    r.FooRequestCount++
}

func (r *RestAPI) barHandler(rw http.ResponseWriter, r *http.Request) {
    r.redisConn.Command("INCR", "barKey")
    r.RequestCount++
    r.BarRequestCount++
}
```

As can be seen, `rest-api` coalesces `http` and `redis` into a simple REST api,
using pre-made library components. `main.go`, the root component, does much the
same:

```go
func main() {
    // Create debug server and start listening in the background
    debugSrv := debug.NewServer()

    // Set up the RestAPI, this will automatically start listening
    restAPI := NewRestAPI()

    // Create another redis connection and use it to store statistics
    statsRedisConn := redis.NewConn("127.0.0.1:6380")
    for {
        time.Sleep(1 * time.Second)
        statsRedisConn.Command("SET", "numReqs", restAPI.RequestCount)
        statsRedisConn.Command("SET", "numFooReqs", restAPI.FooRequestCount)
        statsRedisConn.Command("SET", "numBarReqs", restAPI.BarRequestCount)
    }
}
```

One thing which is clearly missing in this program is proper configuration,
whether from command-line, environment variables, etc.... As it stands, all
configuration parameters, such as the redis addresses and http listen addresses,
are hardcoded. Proper configuration actually ends up being somewhat difficult,
as the ideal case would be for each component to set up the configuration
variables of itself, without its parent needing to be aware. For example,
`redis` could set up `addr` and `pool-size` parameters. The problem is that
there are two `redis` components in the program, and their parameters would
therefore conflict with each other. An elegant solution to this problem is
discussed in the next section.

## Part 2: Context, Configuration, and Runtime

The key to the configuration problem is to recognize that, even if there are two
of the same component in a program, they can't occupy the same place in the
program's structure. In the above example there are two `http` components, one
under `rest-api` and the other under `debug`. Since the structure is represented
as a tree of components, the "path" of any node in the tree uniquely represents
it in the structure. For example, the two `http` components in the previous
example have these paths:

```
root -> rest-api -> http
root -> debug -> http
```

If each component were to know its place in the component tree, then it would
easily be able to ensure that its configuration and initialization didn't
conflict with other components of the same type. If the `http` component sets up
a command-line parameter to know what address to listen on, the two `http`
components in that program would set up:

```
--rest-api-listen-addr
--debug-listen-addr
```

So how can we enable each component to know its path in the component structure?
To answer this we'll have to take a detour through go's `Context` type.

### Context and Configuration

As I mentioned in the Introduction, my example language in this post is Go, but
there's nothing about the concepts I'm presenting which are specific to Go. To
put it simply, Go's builtin `context` package implements a type called
`context.Context` which is, for all intents and purposes, an immutable key/value
store. This means that when you set a key to a value on a Context (using the
`context.WithValue` function) a new Context is returned. The new Context
contains all of the original's key/values, plus the one just set. The original
remains untouched.

(Go's Context also has some behavior built into it surrounding deadlines and
process cancellation, but those aren't relevant for this discussion.)

Context makes sense to use for carrying information about the program's
structure to it's different components; it is informing each of what _context_
it exists in within the larger structure. To use Context effectively, however,
it is necessary to implement some helper functions. Here are their function
signatures:

```go
// NewChild creates and returns a new Context based off of the parent one. The
// child will have a path which is the parent's path appended with the given
// name.
func NewChild(parent context.Context, name string) context.Context

// Path returns the sequence of names which were used to produce this Context
// via calls to the NewChild function.
func Path(ctx context.Context) []string
```

`NewChild` is used to create a new Context, corresponding to a new child node in
the component structure, and `Path` is used retrieve the path of any Context
within that structure. For the sake of keeping the examples simple let's pretend
these functions have been implemented in a package called `mctx`. Here's an
example of how `mctx` might be used in the `redis` component's code:

```go
func NewRedis(ctx context.Context, defaultAddr string) *RedisConn {
    ctx = mctx.NewChild(ctx, "redis")
    ctxPath := mctx.Path(ctx)
    paramPrefix := strings.Join(ctxPath, "-")

    addrParam := flag.String(paramPrefix+"-addr", defaultAddr, "Address of redis instance to connect to")
    // finish setup

    return redisConn
}
```

In our above example, the two `redis` components' parameters would be:

```
// This first parameter is for stats redis, whose parent is the root and
// therefore doesn't have a prefix. Perhaps stats should be broken into its own
// component in order to fix this.
--redis-addr
--rest-api-redis-addr
```

The prefix joining stuff will probably get annoying after a while though, so
let's invent a new package, `mcfg`, which acts like `flag` but is aware of
`mctx`. Then `NewRedis` is reduced to:

```go
func NewRedis(ctx context.Context, defaultAddr string) *RedisConn {
    ctx = mctx.NewChild(ctx, "redis")
    addrParam := flag.String(ctx, "-addr", defaultAddr, "Address of redis instance to connect to")
    // finish setup

    return redisConn
}
```

Sharp-eyed gophers will notice that there's a key piece missing: When is
`mcfg.Parse` called? When does `addrParam` actually get populated? Because you
can't create the redis connection until that happens, but that can't happen
inside `NewRedis` because there might be other things after `NewRedis` which
want to set up parameters. To illustrate the problem, let's look at a simple
program which wants to set up two `redis` components:

```go
func main() {
    // Create the root context, and empty Context.
    ctx := context.Background()

    // Create the Contexts for two sub-components, foo and bar.
    ctxFoo := mctx.NewChild(ctx, "foo")
    ctxBar := mctx.NewChild(ctx, "bar")

    // Now we want to try to create a redis instances for each component. But...

    // This will set up the parameter "--foo-redis-addr", but bar hasn't had a
    // chance to set up its corresponding parameter, so the command-line can't
    // be parsed yet.
    fooRedis := redis.NewRedis(ctxFoo, "127.0.0.1:6379")

    // This will set up the parameter "--bar-redis-addr", but, as mentioned
    // before, NewRedis can't parse command-line.
    barRedis := redis.NewRedis(ctxBar, "127.0.0.1:6379")

    // If the command-line is parsed here, then how can fooRedis and barRedis
    // have been created yet? Creating the redis connection depends on the addr
    // parameters having already been parsed and filled.
}
```

We will solve this problem in the next section.

## Init vs. Start
