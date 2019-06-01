---
title: >-
    Program Structure and Composability
description: >-
    Discussing the nature of program structure, the problems presented by
    complex structures, and a pattern which helps in solving those problems.
---

TODO:
* Double check if I'm using "I" or "We" everywhere (probably should use "I")
* Part 2: Full Example
* Standardize on "programs", not "apps" or "services"
* Prefix all relevant code examples with a package name

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
        globalConns[name] = makeRedisConnection(name)
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
// Package mctx

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
// Package redis

func NewConn(ctx context.Context, defaultAddr string) *RedisConn {
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
`mctx`. Then `redis.NewConn` is reduced to:

```go
// Package redis

func NewConn(ctx context.Context, defaultAddr string) *RedisConn {
    ctx = mctx.NewChild(ctx, "redis")
    addrParam := flag.String(ctx, "-addr", defaultAddr, "Address of redis instance to connect to")
    // finish setup

    return redisConn
}
```

Sharp-eyed gophers will notice that there's a key piece missing: When is
`mcfg.Parse` called? When does `addrParam` actually get populated? Because you
can't create the redis connection until that happens, but that can't happen
inside `redis.NewConn` because there might be other things after `redis.NewConn`
which want to set up parameters. To illustrate the problem, let's look at a
simple program which wants to set up two `redis` components:

```go
func main() {
    // Create the root context, an empty Context.
    ctx := context.Background()

    // Create the Contexts for two sub-components, foo and bar.
    ctxFoo := mctx.NewChild(ctx, "foo")
    ctxBar := mctx.NewChild(ctx, "bar")

    // Now we want to try to create a redis sub-component for each component.

    // This will set up the parameter "--foo-redis-addr", but bar hasn't had a
    // chance to set up its corresponding parameter, so the command-line can't
    // be parsed yet.
    fooRedis := redis.NewConn(ctxFoo, "127.0.0.1:6379")

    // This will set up the parameter "--bar-redis-addr", but, as mentioned
    // before, redis.NewConn can't parse command-line.
    barRedis := redis.NewConn(ctxBar, "127.0.0.1:6379")

    // If the command-line is parsed here, then how can fooRedis and barRedis
    // have been created yet? Creating the redis connection depends on the addr
    // parameters having already been parsed and filled.
}
```

We will solve this problem in the next section.

### Instantiation vs Initialization

Let's break down `redis.NewConn` into two phases: instantiation and initialization.
Instantiation refers to creating the component on the component structure and
having it declare what it needs in order to initialize. After instantiation
nothing external to the program has been done; no IO, no reading of the
command-line, no logging, etc... All that's happened is that the empty shell of
a `redis` component has been created.

Initialization is the phase when that shell is filled. Configuration parameters
are read, startup actions like the creation of database connections are
performed, and logging is output for informational and debugging purposes.

The key to making effective use of this dichotemy is to allow _all_ components
to instantiate themselves before they initialize themselves. By doing this we
can ensure that, for example, all components have had the chance to declare
their configuration parameters before configuration parsing is done.

So let's modify `redis.NewConn` so that it follows this dichotemy. It makes
sense to leave instantiation related code where it is, but we need a mechanism
by which we can declare initialization code before actually calling it. For
this, I will introduce the idea of a "hook".

A hook is, simply a function which will run later. We will declare a new
package, calling it `mrun`, and say that it has two new functions:

```go
// Package mrun

// WithInitHook returns a new Context based off the passed in one, with the //
given hook registered to it.
func WithInitHook(ctx context.Context, hook func()) context.Context

// Init runs all hooks registered using WithInitHook. Hooks are run in the order
// they were registered.
func Init(ctx context.Context)
```

With these two functions we are able to defer the initialization phase of
startup by using the same Contexts we were passing around for the purpose of
denoting component structure. One thing to note is that, since hooks are being
registered onto Contexts within the component instantiation code, the parent
Context will not know about these hooks. Therefore it is necessary to add the
child component's Context back into the parent. To do this we add two final
functions to the `mctx` package:

```go
// Package mctx

// WithChild returns a copy of the parent with the child added to it. Children
// of a Context can be retrieved using the Children function.
func WithChild(parent, child context.Context) context.Context

// Children returns all child Contexts which have been added to the given one
// using WithChild, in the order they were added.
func Children(ctx context.Context) []context.Context
```

Now, with these few extra pieces of functionality in place, let's reconsider the
most recent example, and make a program which creates two redis components which
exist independently of each other:

```go
// Package redis

// NOTE that NewConn has been renamed to WithConn, to reflect that the given
// Context is being returned _with_ a redis component added to it.

func WithConn(parent context.Context, defaultAddr string) (context.Context, *RedisConn) {
    ctx = mctx.NewChild(parent, "redis")

    // we instantiate an empty RedisConn instance and parameters for it. Neither
    // has been initialized yet. They will remain empty until initialization has
    // occurred.
    redisConn := new(RedisConn)
    addrParam := flag.String(ctx, "-addr", defaultAddr, "Address of redis instance to connect to")

    ctx = mrun.WithInitHook(ctx, func() {
        // This hook will run after parameter initialization has happened, and
        // so addrParam will be usable. redisConn will be usable after this hook
        // has run as well.
        *redisConn = makeRedisConnection(*addrParam)
    })

    // Now that ctx has had configuration parameters and intialization hooks
    // instantiated into it, return both it and the empty redisConn instance
    // back to the parent.
    return mctx.WithChild(parent, ctx), redisConn
}

////////////////////////////////////////////////////////////////////////////////

// Package main

func main() {
    // Create the root context, an empty Context.
    ctx := context.Background()

    // Create the Contexts for two sub-components, foo and bar.
    ctxFoo := mctx.NewChild(ctx, "foo")
    ctxBar := mctx.NewChild(ctx, "bar")

    // Add redis components to each of the foo and bar sub-components. The
    // returned Contexts will be used to initialize the redis components.
    ctxFoo, redisFoo := redis.WithConn(ctxFoo, "127.0.0.1:6379")
    ctxBar, redisBar := redis.WithConn(ctxBar, "127.0.0.1:6379")

    // Add the sub-component contexts back to the root, so they can all be
    // initialized at once.
    ctx = mctx.WithChild(ctx, ctxFoo)
    ctx = mctx.WithChild(ctx, ctxBar)

    // Parse will descend into the Context and all of its children, discovering
    // all registered configuration parameters and filling them from the
    // command-line.
    mcfg.Parse(ctx)

    // Now that configuration has been initialized, run the Init hooks for each
    // of the sub-components.
    mrun.Init(ctx)

    // At this point the redis components have been fully initialized and may be
    // used. For this example we'll copy all keys from one to the other.
    keys := redisFoo.Command("KEYS", "*")
    for i := range keys {
        val := redisFoo.Command("GET", keys[i])
        redisBar.Command("SET", keys[i], val)
    }
}
```

### Full example

## Part 3: Annotations, Logging, and Errors

Let's shift gears away from the component structure for a bit, and talk about a
separate, but related, set of issues: those related to logging and errors.

Both logging and error creation share the same problem, that of collecting as
much contextual information around an event as possible. This is often done
through string formatting, like so:

```go
// ServeHTTP implements the http.Handler method and is used to serve App's HTTP
// endpoints.
func (app *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    log.Printf("incoming request from remoteAddr:%s for url:%s", r.RemoteAddr, r.URL.String())

    // begin actual request handling
}
```

In this example the code is logging an event, an incoming HTTP request, and
including contextual information in that log about the remote address of the
requester and the URL being requested.

Similarly, an error might be created like this:

```go
func (app *App) GetUsername(userID int) (string, error) {
    userName, err := app.Redis.Command("GET", userID)
    if err != nil {
        return "", fmt.Errorf("could not get username for userID:%d: %s", userID, err)
    }
    return userName, nil
}
```

In that example, when redis returns an error the error is extended to include
contextual information about what was attempting to be done (`could not get
username`) and the userID involved. In newer versions of Go, and indeed in many
other programming languages, the error will also include information about where
in the source code it occurred, such as file name and line number.

It is my experience that both logging and error creation often take up an
inordinate amount of space in many programs. This is due to a desire to
contextualize as much as possible, since in a large program it can be difficult
to tell exactly where something is happening, even if you're looking at the log
entry or error. For example, if a program has a set of HTTP endpoints, each one
performing a redis call, what good is it to see the log entry `redis command had
an error: took too long` without also knowing which command is involved, and
which endpoint is calling it? Very little.

So many programs of this nature end up looking like this:

```go
func (app *App) httpEndpointA(rw http.ResponseWriter, r *http.Request) {
    err := app.Redis.Command("SET", "foo", "bar")
    if err != nil {
        log.Printf("redis error occurred in EndpointA, calling SET: %s", err)
    }
}

func (app *App) httpEndpointB(rw http.ResponseWriter, r *http.Request) {
    err := app.Redis.Command("INCR", "baz")
    if err != nil {
        log.Printf("redis error occurred in EndpointA, calling INCR: %s", err)
    }
}

// etc...
```

Obviously logging is taking up the majority of the code-space in those examples,
and that doesn't even include potentially pertinent information such as IP
address.

Another aspect of the logging/error dichotemy is that they are often dealing in
essentially the same data. This makes sense, as both are really dealing with the
same thing: capturing context for the purpose of later debugging. So rather than
formatting strings by hand for each use-case, let's instead use our friend,
`context.Context`, to carry the data for us.

### Annotations

I will here introduce the idea of "annotations", which are essentially key/value
pairs which can be attached to a Context and retrieved later. To implement
annotations I will introduce two new functions to the `mctx` package:

```go
// Package mctx

// Annotate returns a new Context with the given key/value pairs embedded into
// it, which can be later retrieved using the Annotations method. If any keys
// conflict with previous annotations, their values will overwrite the
// previously annotated values for those keys.
func Annotate(ctx context.Context, keyvals ...interface{}) context.Context

// Annotations returns all annotations which have been set on the Context using
// Annotate.
func Annotations(ctx context.Context) map[interface{}]interface{}
```
