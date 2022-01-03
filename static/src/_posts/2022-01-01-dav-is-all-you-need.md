---
title: >-
    DAV is All You Need
description: >-
    Contacts, calendars, passwords, oh my!
tags: tech
---

For some time now I've been trying to find an alternative solution to Google
Keep for shared note taking. The motivation for this change was two-fold:

* Google sucks, and I'm trying to get their products out of my life in favor of
  self-hosted options.

* Google Keep _really_ sucks. Seriously, it can barely load on my Chromebook
  because of whatever bloated ass web framework they're using for it. It's just
  a note taking app!

So this weekend I buckled down and actually made the switch. The first step was
to find something to switch _to_, however, which ended up being not trivial.
There's a million different options in this space, but surprisingly few which
could fulfill the exact niche we need in our household:

* Fully open-source and open protocol. If it's not open it's not worth the
  bother of switching, cause we'll just have to do it all again once whatever
  product we switch to gets acqui-hired by a food delivery app.

* Self-hosted using a _simple_ server-side component. I'm talking something that
  listens on a public port and saves data to a file on disk, and _that's it_.
  No database processes, no message queues, no bullshit. We're not serving a
  million users here, there's no reason to broaden the attack surface
  unnecessarily.

* Multi-platform support, including mobile. Our primary use-case here is our
  grocery list, which needs to be accessible by everyone everywhere.

I've already got a Nextcloud instance running at home, and there is certainly a
notes extension for it, so that could have been an option here. But Nextcloud
very much does not fall into the second point above: it's not simple. It's a
giant PHP app that uses Postgres as a backend, has its own authentication and
session system, and has a plugin system. Frankly, it was easily the biggest
security hole on the entire server, and I wasn't eager to add usage to it.

Happily, I found another solution.

## WebDAV

There's a project called [Joplin](https://joplinapp.org/) which implements a
markdown-based notes system with clients for Android, iPhone, Linux, Mac, and
Windows. Somewhat interestingly there is _not_ a web client for it, but on
further reflection I don't think that's a big deal... no bloated javascript
frameworks to worry about at least.

In addition to their own cloud backend, Joplin supports a number of others, with
the most interesting being WebDAV. WebDAV is an XML-based extension to HTTP
which allows for basic write operations on the server-side, and which uses
HTTP's basic auth for authentication. You can interact with it using curl if you
like, it really can't get simpler.

[Caddy](https://caddyserver.com/) is the server I use to handle all incoming
HTTP requests to my server, and luckily there's a semi-official
[WebDAV](https://github.com/mholt/caddy-webdav) plugin which adds WebDAV
support. With that compiled in, the `Caddyfile` configuration is nothing more
than:

```
hostname.com {

    route {

        basicauth {
            sharedUser sharedPassword
        }


        webdav {
            root /data/webdav
        }

    }

}
```

With that in place, any Joplin client can be pointed at `hostname.com` using the
shared username/assword, and all data is stored directly to `/data/webdav` by
Caddy. Easy-peasy.

## CardDAV/CalDAV

Where WebDAV is an extension of HTTP to allow for remotely modifying files
genearlly, CardDAV and CalDAV are extensions of WebDAV for managing remote
stores of contacts and calendar events, respectively. At least, that's my
understanding.

Nextcloud has its own Web/Card/CalDAV service, and that's what I had been, up
till this point, using for syncing my contacts and calendar from my phone. But
now that I was setting up a separate WebDAV endpoint, I figured it'd be worth
setting up a separate Card/CalDAV service and get that much closer to getting
off Nextcloud entirely.

There is, as far as I know, no Card or CalDAV extension for Caddy, so I'd still
need a new service running. I came across
[radicale](https://radicale.org/v3.html), which fits the bill nicely. It's a
simple CalDAV and CardDAV server which saves directly to disk, much like the
Caddy WebDAV plugin. With that running, I needed only to add the following to my
`Caddyfile`, above the `webdav` directive:

```
handle /radicale/* {

    uri strip_prefix /radicale

    reverse_proxy 127.0.0.1:5454 {
        header_up X-Script-Name /radicale
    }

}
```

Now I could point the [DAVx5](https://www.davx5.com/) app on my phone to
`hostname.com/radicale` and boom, contact and calendar syncing was within reach.
I _did_ have a lot of problems getting DAVx5 working properly, but those were
more to do with Android than self-hosting, and I eventually worked through them.

## Passwords

At this point I considered that the only thing I was still really using
Nextcloud for was password management, a la Lastpass or 1Password. I have a lot
of gripes with Nextcloud's password manager, in addition to my aforementioned
grips with Nextcloud generally, so I thought it was worth seeing if some DAV or
another could be the final nail in Nextcloud's coffin.

A bit of searching around led me to [Tusk](https://subdavis.com/Tusk/), a chrome
extension which allows the chrome browser to fetch a
[KeePassXC](https://keepassxc.org/) database from a WebDAV server, decode it,
and autofill it into a website. Basically perfect. I had only to export my
passwords from Nextcloud as a CSV, import them into a fresh KDBX file using the
KeePassXC GUI, place the file in my WebDAV folder, and point Tusk at that.

I found the whole experience of using Tusk to be extremely pleasant. Everything
is very well labeled and described, and there's appropriate warnings and such in
places where someone might commit a security crime (e.g. using the same password
for WebDAV and their KDBX file).

My one gripe is that it seems to be very slow to unlock the file in practice. I
don't _think_ this has to do with my server, as Joplin is quite responsive, so
it could instead have to do with my KDBX file's decryption difficulty setting.
Perhaps Tusk is doing the decryption in userspace javascript... I'll have to
play with it some.

But it's a small price to be able to turn off Nextcloud completely, which I have
now done. I can sleep easier at night now, knowing there's not some PHP
equivalent to Log4j which is going to bite me in the ass one day while I'm on
vacation.
