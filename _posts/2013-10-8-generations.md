---
layout: post
title: Generations
categories:
- crypticio
---

A simple file distribution strategy for very large scale, high-availability
file-services.

# The problem

Working at a shop where we have millions of different files, any of which could
be arbitrarily chosen to serve to a file at any given time. These files are
uploaded by users of the app and retrieved by others.

Scaling such a system is no easy task. The chosen solution involves shuffling
files around on a nearly constant basis, making sure that files which are more
"popular" are on fast drives, while at the same time making sure that no drives
are at capicty and at the same time that all files, even newly uploaded ones,
are stored redundantly.

The problem with this solution is one of coordination. At any given moment the
app needs to be able to "find" a file so it can give the client a link to
download the file from one of the servers that it's on. Full-filling this simple
requirement means that all datastores/caches where information about where a
file lives need to be up-to-date at all times, and even then there are
race-conditions and network failures to contend with, while at all times the
requirements of the app evolve and change.

# A simpler solution

Let's say you want all files which get uploaded to be replicated in triplicate
in some capacity. You buy three identical hard-disks, and put each on a separate
server.  As files get uploaded by clients, each file gets put on each drive
immediately. When the drives are filled (which should be at around the same
time), you stop uploading to them.

That was generation 0.

You buy three more drives, and start putting all files on them instead. This is
going to be generation 1. Repeat until you run out of money.

That's it.

## That's it?

It seems simple and obvious, and maybe it's the standard thing which is done,
but as far as I can tell no-one has written about it (though I'm probably not
searching for the right thing, let me know if this is the case!).

## Advantages

* It's so simple to implement, you could probably do it in a day if you're
starting a project from scratch

* By definition of the scheme all files are replicated in multiple places.

* Minimal information about where a file "is" needs to be stored. When a file is
uploaded all that's needed is to know what generation it is in, and then what
nodes/drives are in that generation.

* Drives don't need to "know" about each other. What I mean by this is that
whatever is running as the receive point for file-uploads on each drive doesn't
have to coordinate with its siblings running on the other drives in the
generation. In fact it doesn't need to coordinate with anyone. You could
literally rsync files onto your drives if you wanted to. I would recommend using
[marlin][0] though :)

* Scaling is easy. When you run out of space you can simply start a new
generation. If you don't like playing that close to the chest there's nothing to
say you can't have two generations active at the same time.

* Upgrading is easy. As long as a generation is not marked-for-upload, you can
easily copy all files in the generation into a new set of bigger, badder drives,
add those drives into the generation in your code, remove the old ones, then
mark the generation as uploadable again.

* Distribution is easy. You just copy a generation's files onto a new drive in
Europe or wherever you're getting an uptick in traffic from and you're good to
go.

* Management is easy. It's trivial to find out how many times a file has been
replicated, or how many countries it's in, or what hardware it's being served
from (given you have easy access to information about specific drives).

## Caveats

The big caveat here is that this is just an idea. It has NOT been tested in
production. But we have enough faith in it that we're going to give it a shot at
cryptic.io. I'll keep this page updated.

The second caveat is that this scheme does not inherently support caching. If a
file suddenly becomes super popular the world over your hard-disks might not be
able to keep up, and it's probably not feasible to have an FIO drive in *every*
generation. I think that [groupcache][1] may be the answer to this problem,
assuming your files are reasonably small, but again I haven't tested it yet.

[0]: https://github.com/cryptic-io/marlin
[1]: https://github.com/golang/groupcache
