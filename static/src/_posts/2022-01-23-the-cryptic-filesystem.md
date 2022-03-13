---
title: >-
    The Cryptic Filesystem
description: >-
    Hey, I'm brainstorming here!
series: nebula
tags: tech
---

Presently the cryptic-net project has two components: a VPN layer (implemented
using [nebula][nebula], and DNS component which makes communicating across that
VPN a bit nicer. All of this is wrapped up in a nice bow using an AppImage and a
simple process manager. The foundation is laid for adding the next major
component: a filesystem layer.

I've done a lot of research and talking about this layer, and you can see past
posts in this series talking about it. Unfortunately, I haven't really made much
progress on a solution. It really feels like there's nothing out there already
implemented, and we're going to have to do it from scratch.

To briefly recap the general requirements of the cryptic network filesystem
(cryptic-fs), it must have:

* Sharding of the fs dataset, so each node doesn't need to persist the full
  dataset.

* Replication factor (RF), so each piece of content must be persisted by at
  least N nodes of the clusters.

* Nodes are expected to be semi-permanent. They are expected to be in it for the
  long-haul, but they also may flit in and out of existence frequently.

* Each cryptic-fs process should be able to track multiple independent
  filesystems, with each node in the cluster not necessarily tracking the same
  set of filesystems as the others.

This post is going to be a very high-level design document for what, in my head,
is the ideal implementation of cryptic-fs. _If_ cryptic-fs is ever actually
implemented it will very likely differ from this document in major ways, but one
must start somewhere.

[nebula]: https://github.com/slackhq/nebula

## Merkle DAG

It wouldn't be a modern network filesystem project if there wasn't a [Merkle
DAG][mdag]. The minutia of how a Merkle DAG works isn't super important here,
the important bits are:

* Each file is represented by a content identifier (CID), which is essentially a
  consistent hash of the file's contents.

* Each directory is also represented by a CID which is generated by hashing the
  CIDs of the directory's files and their metadata.

* Since the root of the filesystem is itself a directory, the entire filesystem
  can be represented by a single CID. By tracking the changing root CID all
  hosts participating in the network filesystem can cheaply identify the latest
  state of the entire filesystem.

A storage system for a Merkle DAG is implemented as a key-value store which maps
CID to directory node or file contents. When nodes in the cluster communicate
about data in the filesystem they will do so using these CIDs; one node might
ask the other "can you give me CID `AAA`", and the other would respond with the
contents of `AAA` without really caring about whether or not that CID points to
a file or directory node or whatever. It's quite a simple system.

As far as actual implementation of the storage component, it's very likely we
could re-use some part of the IPFS code-base rather than implementing this from
scratch.

[mdag]: https://docs.ipfs.io/concepts/merkle-dag/

## Consensus

The cluster of nodes needs to (roughly) agree on some things in order to
function:

* What the current root CID of the filesystem is.

* Which nodes have which CIDs persisted.

These are all things which can change rapidly, and which _every_ node in the
cluster will need to stay up-to-date on. On the other hand, given efficient use
of the boolean tagged CIDs mentioned in the previous section, this is a dataset
which could easily fit in memory even for large filesystems.

I've done a bunch of research here and I'm having trouble finding anything
existing which fits the bill. Most databases expect the set of nodes to be
pretty constant, so that eliminates most of them. Here's a couple of other ideas
I spitballed:

* Taking advantage of the already written [go-ds-crdt][crdt] package which the
  [IPFS Cluster][ipfscluster] project uses. My biggest concern with this
  project, however, is that the entire history of the CRDT must be stored on
  each node, which in our use-case could be a very long history.

* Just saying fuck it and using a giant redis replica-set, where each node in
  the cluster is a replica and one node is chosen to be the primary. [Redis
  sentinel][sentinel] could be used to decide the current primary. The issue is
  that I don't think sentinel is designed to handle hundreds or thousands of
  nodes, which places a ceiling on cluster capacity. I'm also not confident that
  the primary node could handle hundreds/thousands of replicas syncing from it
  nicely; that's not something Redis likes to do.

* Using a blockchain engine like [Tendermint][tendermint] to implement a custom,
  private blockchain for the cluster. This could work performance-wise, but I
  think it would suffer from the same issue as CRDT.

It seems to me like some kind of WAN-optimized gossip protocol would be the
solution here. Each node already knows which CIDs it itself has persisted, so
what's left is for all nodes to agree on the latest root CID, and to coordinate
who is going to store what long-term.

[crdt]: https://github.com/ipfs/go-ds-crdt
[ipfscluster]: https://cluster.ipfs.io/
[sentinel]: https://redis.io/topics/sentinel
[tendermint]: https://tendermint.com/

### Gossip

The [gossipsub][gossipsub] library which is built into libp2p seems like a good
starting place. It's optimized for WANs and, crucially, is already implemented.

Gossipsub makes use of different topics, onto which peers in the cluster can
publish messages which other peers who are subscribed to those topics will
receive. It makes sense to have a topic-per-filesystem (remember, from the
original requirements, that there can be multiple filesystems being tracked), so
that each node in the cluster can choose for itself which filesystems it cares
to track.

The messages which can get published will be dependent on the different
situations in which nodes will want to communicate, so it's worth enumerating
those.

**Situation #1: Node A wants to obtain a CID**: Node A will send out a
`WHO_HAS:<CID>` message (not the actual syntax) to the topic. Node B (and
possibly others), which has the CID persisted, will respond with `I_HAVE:<CID>`.
The response will be sent directly from B to A, not broadcast over the topic,
since only A cares. The timing of B's response to A could be subject to a delay
based on B's current load, such that another less loaded node might get its
response in first.

From here node A would initiate a download of the CID from B via a direct
connection. If node A has enough space then it will persist the contents of the
CID for the future.

This situation could arise because the user has opened a file in the filesystem
for reading, or has attempted to enumerate the contents of a directory, and the
local storage doesn't already contain that CID.

**Situation #2: Node A wants to delete a CID which it has persisted**: Similar
to #1, Node A needs to first ensure that other nodes have the CID persisted, in
order to maintain the RF across the filesystem. So node A first sends out a
`WHO_HAS:<CID>` message. If >=RF nodes respond with `I_HAVE:<CID>` then node A
can delete the CID from its storage without concern. Otherwise it should not
delete the CID.

**Situation #2a: Node A wants to delete a CID which it has persisted, and which
is not part of the current filesystem**: If the filesystem is in a state where
the CID in question is no longer present in the system, then node A doesn't need
to care about the RF and therefore doesn't need to send any messages.

**Situation #3: Node A wants to update the filesystem root CID**: This is as
simple as sending out a `ROOT:<CID>` message on the topic. Other nodes will
receive this and note the new root.

**Situation #4: Node A wants to know the current filesystem root CID**: Node A
sends out a `ROOT?` message. Other nodes will respond to node A directly telling
it the current root CID.

These describe the circumstances around the messages used across the gossip
protocol in a very shallow way. In order to properly flesh out the behavior of
the consistency mechanism we need to dive in a bit more.

### Optimizations, Replication, and GC

A key optimization worth hitting straight away is to declare that each node will
always immediately persist all directory CIDs whenever a `ROOT:<CID>` message is
received. This will _generally_ only involve a couple of round-trips with the
host which issued the `ROOT:<CID>` message, with opportunity for
parallelization.

This could be a problem if the directory structure becomes _huge_, at which
point it might be worth placing some kind of limit on what percent of storage is
allowed for directory nodes. But really... just have less directories people!

The next thing to dive in on is replication. We've already covered in situation
 #1 what happens if a user specifically requests a file. But that's not enough
to ensure the RF of the entire filesystem, as some files might not be requested
by any users except the original user to add the file.

We can note that each node knows when a file has been added to the filesystem,
thanks to each node knowing the full directory tree. So upon seeing that a new
file has been added, a node can issue a `WHO_HAS:<CID>` message for it, and if
less than RF nodes respond then it can persist the CID. This is all assuming
that the node has enough space for the new file.

One wrinkle in that plan is that we don't want all nodes to send the
`WHO_HAS:<CID>` at the same time for the same CID, otherwise they'll all end up
downloading the CID and over-replicating it. A solution here is for each node to
delay it's `WHO_HAS:<CID>` based on how much space it has left for storage, so
nodes with more free space are more eager to pull in new files.

Additionally, we want to have nodes periodically check the replication status of
each CID in the filesystem. This is because nodes might pop in and out of
existence randomly, and the cluster needs to account for that. The way this can
work is that each node periodically picks a CID at random and checks the
replication status of it. If the period between checks is calculated as being
based on number of online nodes in the cluster and the number of CIDs which can
be checked, then it can be assured that all CIDs will be checked within a
reasonable amount of time with minimal overhead.

This dovetails nicely with garbage collection. Given that nodes can flit in and
out of existence, a node might come back from having been down for a time, and
all CIDs it had persisted would then be over-replicated. So the same process
which is checking for under-replicated files will also be checking for
over-replicated files.

### Limitations

This consistency mechanism has a lot of nice properties: it's eventually
consistent, it nicely handles nodes coming in and out of existence without any
coordination between the nodes, and it _should_ be pretty fast for most cases.
However, it has its downsides.

There's definitely room for inconsistency between each node's view of the
filesystem, especially when it comes to the `ROOT:<CID>` messages. If two nodes
issue `ROOT:<CID>` messages at the same time then it's extremely likely nodes
will have a split view of the filesystem, and there's not a great way to
resolve this until another change is made on another node. This is probably the
weakest point of the whole design.

[gossipsub]: https://github.com/libp2p/specs/tree/master/pubsub/gossipsub

## FUSE

The final piece is the FUSE connector for the filesystem, which is how users
actually interact with each filesystem being tracked by their node. This is
actually the easiest component, if we use an idea borrowed from
[Tahoe-LAFS][tahoe], cryptic-fs can expose an SFTP endpoint and that's it.

The idea is that hooking up an existing SFTP implementation to the rest of
cryptic-fs should be pretty straightforward, and then every OS should already
have some kind of mount-SFTP-as-FUSE mechanism already, either built into it or
as an existing application. Exposing an SFTP endpoint also allows a user to
access the cryptic-fs remotely if they want to.

[tahoe]: https://tahoe-lafs.org/trac/tahoe-lafs

## Ok

So all that said, clearly the hard part is the consistency mechanism. It's not
even fully developed in this document, but it's almost there. The next step,
beyond polishing up the consistency mechanism, is going to be roughly figuring
out all the interfaces and types involved in the implementation, planning out
how those will all interact with each other, and then finally an actual
implementation!