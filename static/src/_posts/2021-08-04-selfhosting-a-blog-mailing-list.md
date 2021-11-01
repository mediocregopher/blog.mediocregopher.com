---
title: >-
    Self-Hosting a Blog Mailing List
description: >-
    For fun and no profit.
tags: tech
---

As of this week the Mediocre Blog has a new follow mechanism: email! [Sign up
on the **Follow** page][follow] and you'll get an email everytime a new post
is published to the blog. It's like RSS, except there's a slight chance you
might actually use it.

This post will detail my relatively simple setup for this, linking to points
within my blog's server code which are relevant. While I didn't deliberately
package my code up into a nice public package, if you know have some cursory
knowledge of Go you could probably rip my code and make it work for you. Don't
worry, it has a [permissive license](/assets/wtfpl.txt).

[follow]: /follow.html

## Email Server

Self-hosting email is the hardest and most foreign part of this whole
thing for most devs. The long and the short of it is that it's very unlikely you
can do this without renting a VPS somewhere. Luckily there are VPSs out there
which are cheap and which allow SMTP traffic, so it's really just a matter of
biting the cost bullet and letting your definition of "self-hosted" be a bit
flexible. At least you still control the code!

I highly recommend [maddy][maddy] as an email server which has everything you
need out of the box, no docker requirements, and a flexible-yet-simple
configuration language.  I've discussed [in previous posts][maddypost] the
general steps I've used to set up maddy on a remote VPS, and so I won't
re-iterate here. Just know that I have a VPS on my private [nebula][nebula] VPN,
with a maddy server listening for outgoing mail on port 587, with
username/password authentication on that port.

[maddy]: https://maddy.email
[maddypost]: {% post_url 2021-07-06-maddy-vps %}
[nebula]: https://github.com/slackhq/nebula

## General API Design

The rest of the system lies within the Go server which hosts my blog. There is
only a single instance of the server, and it runs in my living room. With these
as the baseline environmental requirements, the rest of the design follows
easily:

* The Go server provides [three REST API endpoints][restendpoints]:

    - `POST /api/mailinglist/subscribe`: Accepts a POST form argument `email`, sends a
      verification email to that email address.

    - `POST /api/mailinglist/finalize`: Accepts a POST form argument `subToken`,
      which is a random token sent to the user when they subscribe. Only by
      finalizing their subscription can a user be considered actually
      subscribed.

    - `POST /api/mailinglist/unsubscribe`: Accepts a POST form argument
      `unsubToken`, which is sent with each blog post notification to the user.

* The static frontend code has [two pages][staticpages] related to the mailing
  list:

    - `/mailinglist/finalize.html`: The verification email which is sent to the
      user links to this page, with the `subToken` as a GET argument. This page
      then submits the `subToken` to the `POST /api/mailinglist/finalize`
      endpoint.

    - `/mailinglist/unsubscribe.html`: Each blog post notification email sent to
      users contains a link to this page, with an `unsubToken` as a GET
      argument. This page then submits the `unsubToken` to the `POST
      /api/mailinglist/unsubscribe` endpoint.

It's a pretty small API, but it covers all the important things, namely
verification (because I don't want people signed up against their will, nor do I
want to be sending emails to fake email addresses), and unsubscribing.

[restendpoints]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/5ca7dadd02fb49dd62ad448d12021359e41beec1/srv/cmd/mediocre-blog/main.go#L169
[staticpages]: https://github.com/mediocregopher/blog.mediocregopher.com/tree/9c3ea8dd803d6f0df768e3ae37f8c4ab2efbcc5c/static/src/mailinglist

## Proof-of-work

It was important to me that someone couldn't just sit and submit random emails
to the `POST /api/mailinglist/subscribe` endpoint in a loop, causing my email
server to eventually get blacklisted. To prevent this I've implemented a simple
proof-of-work (PoW) system, whereby the client must first obtain a PoW
challenge, generate a solution for that challenge (which involves a lot of CPU
time), and then submit that solution as part of the subscribe endpoint call.

Both the [server-side][powserver] and [client-side][powclient] code can be found
in the blog's git repo. You could theoretically view the Go documentation for
the server code on pkg.go.dev, but apparently their bot doesn't like my WTFPL.

When providing a challenge to the client, the server sends back two values: the
seed and the target.

The target is simply a number whose purpose will become apparent in a second.

The seed is a byte-string which encodes:

* Some random bytes.

* An expiration timestamp.

* A target (matching the one returned to the client alongside the seed).

* An HMAC-MD5 which signs all of the above.

When the client submits a valid solution the server checks the HMAC to ensure
that the seed was generated by the server, it checks the expiration to make sure
the client didn't take too long to solve it, and it checks in an [internal
storage][powserverstore] whether that seed hasn't already been solved. Because
the expiration is built into the seed the server doesn't have to store each
solved seed forever, only until the seed has expired.

To generate a solution to the challenge the client does the following:

* Concatenate up to `len(seed)` random bytes onto the original seed given by the
  server.

* Calculate the SHA512 of that.

* Parse the first 4 bytes of the resulting hash as a big-endian uint32.

* If that uint32 is less than the target then the random bytes generated in the
  first step are a valid solution. Otherwise the client loops back to the first
  step.

Finally, a new endpoint was added: `GET /api/pow/challenge`, which returns a PoW
seed and target for the client to solve. Since seeds don't require storage in a
database until _after_ they are solved there are essentially no consequences to
someone spamming this in a loop.

With all of that in place, the `POST /api/mailinglist/subscribe` endpoint
described before now also requires a `powSeed` and a `powSolution` argument. The
[Follow][follow] page, prior to submitting a subscribe request, first retrieves
a PoW challenge, generates a solution, and only _then_ will it submit the
subscribe request.

[powserver]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/9c3ea8dd803d6f0df768e3ae37f8c4ab2efbcc5c/srv/pow/pow.go
[powserverstore]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/5ca7dadd02fb49dd62ad448d12021359e41beec1/srv/pow/store.go
[powclient]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/9c3ea8dd803d6f0df768e3ae37f8c4ab2efbcc5c/static/src/assets/solvePow.js

## Storage

Storage of emails is fairly straightforward: since I'm not running this server
on multiple hosts, I can just use [SQLite][sqlite]. My code for storage in
SQLite can all be found [here][sqlitecode].

My SQLite table has a single table:

```
CREATE TABLE emails (
	id          TEXT PRIMARY KEY,
	email       TEXT NOT NULL,
	sub_token   TEXT NOT NULL,
	created_at  INTEGER NOT NULL,

	unsub_token TEXT,
	verified_at INTEGER
)
```

It will probably one day need an index on `sub_token` and `unsub_token`, but I'm
not quite there yet.

The `id` field is generated by first lowercasing the email (because emails are
case-insensitive) and then hashing it. This way I can be sure to identify
duplicates easily. It's still possible for someone to do the `+` trick to get
their email in multiple times, but as long as they verify each one I don't
really care.

[sqlite]: https://sqlite.org/index.html
[sqlitecode]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/5ca7dadd02fb49dd62ad448d12021359e41beec1/srv/mailinglist/store.go

## Publishing

Publishing is quite easy: my [MailingList interface][mailinglistinterface] has a
`Publish` method on it, which loops through all records in the SQLite table,
discards those which aren't verified, and sends an email to the rest containing:

* A pleasant greeting.

* The new post's title and URL.

* An unsubscribe link.

I will then use a command-line interface to call this `Publish` method. I
haven't actually made that interface yet, but no one is subscribed yet so it
doesn't matter.

[mailinglistinterface]: https://github.com/mediocregopher/blog.mediocregopher.com/blob/5ca7dadd02fb49dd62ad448d12021359e41beec1/srv/mailinglist/mailinglist.go#L23

## Easy-Peasy

The hardest part of the whole thing was probably getting maddy set up, with a
close second being trying to decode a hex string to a byte string in javascript
(I tried Crypto-JS, but it wasn't working without dragging in webpack or a bunch
of other nonsense, and vanilla JS doesn't have any way to do it!).

Hopefully reading this will make you consider self-hosting your own blog's
mailing list as well. If we let these big companies keep taking over all
internet functionality then eventually they'll finagle the standards so that
no one can self-host anything, and we'll have to start all over.

And really, do you _need_ tracking code on the emails you send out for your
recipe blog? Just let your users ignore you in peace and quiet.
