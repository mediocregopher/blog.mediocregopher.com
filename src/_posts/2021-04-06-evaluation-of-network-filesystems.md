---
title: >-
    Evaluation of Network Filesystems
description: >-
    There can only be one.
series: nebula
tags: tech
---

It's been a bit since updating my progress on what I've been lately calling the
"cryptic nebula" project. When I last left off I was working on building the
[mobile nebula][mobile_nebula] using [nix][nix]. For the moment I gave up on
that dream, as flutter and nix just _really_ don't get along and I don't want to
get to distracted on problems that aren't critical to the actual goal.

Instead I'd like to pursue the next critical component of the system, and
that's a shared filesystem. The use-case I'm ultimately trying to achieve is:

* All hosts communicate with each other via the nebula network.
* All hosts are personal machines owned by individuals, _not_ cloud VMs.
* A handful of hosts are always-on, or at least as always-on as can be achieved
  in a home environment.
* All hosts are able to read/write to a shared filesystem, which is mounted via
  FUSE (or some other mechanism, though I can't imagine what) on their computer.
* Top-level directories within the shared filesystem can be restricted, so
  that only a certain person (or host) can read/write to them.

What I'm looking for is some kind of network filesystem, of which there are
_many_. This document will attempt to evaluate all relevant projects and come up
with the next steps. It may be that no project fits the bill perfectly, and that
I'm stuck either modifying an existing project to my needs or, if things are
looking really dire, starting a new project.

The ultimate use-case here is something like a self-hosted, distributed [keybase
filesystem](https://book.keybase.io/docs/files); somewhere where individuals in
the cluster can back up their personal projects, share files with each other,
and possibly even be used as the base layer for more complex applications on
top.

The individuals involved shouldn't have to deal with configuring their
distributed FS, either to read from it or add storage resources to it. Ideally
the FS process can be bundled together with the nebula process and run opaquely;
the user is just running their "cryptic nebula" process and everything else is
handled in the background.

## Low Pass Filter

There are some criteria for these projects that I'm not willing to compromise
on; these criteria will form a low pass filter which, hopefully, will narrow our
search appreciably.

The network filesystem used by the cryptic nebula must:

* Be able to operate over a nebula network (obviously).
* Be open-source. The license doesn't matter, as long as the code is available.
* Run on both Mac and Linux.
* Not require a third-party to function.
* Allows for a replication factor of 3.
* Supports sharding of data (ie each host need not have the entire dataset).
* Allow for mounting a FUSE filesystem in any hosts' machine to interact with
  the network filesystem.
* Not run in the JVM, or any other VM which is memory-greedy.

The last may come across as mean, but the reason for it is that I forsee the
network filesystem client running on users' personal laptops, which cannot be
assumed to have resources to spare.

## Rubric

Each criteria in the next set lies along a spectrum. Any project may meet one of
thses criteria fully, partially, or not at all. For each criteria I assign a
point value according to how fully a project meets the criteria, and then sum up
the points to give the project a final score. The project with the highest final
score is not necessarily the winner, but this system should at least give some
good candidates for final consideration.

The criteria, and their associated points values, are:

* **Hackability**: is the source-code of the project approachable?
    - 0: No
    - 1: Kind of, and there's not much of a community.
    - 2: Kind of, but there is an active community.
    - 3: Yes

* **Documentation**: is the project well documented?
    - 0: No docs.
    - 1: Incomplete or out-of-date docs.
    - 2: Very well documented.

* **Transience**: how does the system handle hosts appearing or disappearing?
    - 0: Requires an automated system to be built to handle adding/removing
      hosts.
    - 1: Gracefully handled.

* **Priority**: is it possible to give certain hosts priority when choosing
  which will host/replicate some piece of data?
    - 0: No.
    - 1: Yes.

* **Caching**: will hosts reading a file have that file cached locally for the
  next reading (until the file is modified)?
    - 0: No.
    - 1: Yes.

* **Conflicts**: if two hosts updated the same file at the same time, how is
  that handled?
    - 0: The file can no longer be updated.
    - 1: One update clobbers the other, or both go through in an undefined
      order.
    - 2: One update is disallowed.
    - 3: A copy of the file containing the "losing" update is created (ie: how
      dropbox does it).
    - 4: Strategy can be configured on the file/directory level.

* **Consistency**: how does the system handle a file being changed frequently?
    - 0: File changes must be propagated before subsequent updates are allowed (fully consistent).
    - 1: Files are snapshotted at some large-ish interval (eventually consistent).
    - 2: File state (ie content hash, last modifid, etc) is propagated
      frequently but contents are only fully propagated once the file has
      "settled" (eventually consistent with debounce).

* **POSIX**: how POSIX compliant is the mounted fileystem?
    - 0: Only the most basic features are implemented.
    - 1: Some extra features are implemented.
    - 2: Fully POSIX compliant.

* **Scale**: how many hosts can be a part of the cluster?
    - 0: A finite number.
    - 1: A finite number of dedicated hosts, infinite ephemeral.
    - 2: Infinite hosts.

* **Failure**: how does the system handle failures (network partitions, hosts
  hanging, buggy client versions)?
    - 0: Data loss.
    - 1: Reads and writes are halted.
    - 2: Reads are allowed but writes are halted.
    - 3: System is partially read/write, except effected parts.

* **Limitations**: are there limits on how big files can be, or how big
  directories can be?
    - 0: Files are limited to below 1TB in size.
    - 1: Directories are limited to below 100,000 files.
    - 2: No limits.

* **Encryption**: how is data encrypted?
    - 0: Not at all, DIY.
    - 1: Encrypted at rest.
    - 2: Per-user encryption.

* **Permissions**: how are modifications to data restricted?
    - 0: Not at all.
    - 1: Permissions are only superifically enforced.
    - 2: Fully enforced user/group restrictions, complex patterns, and/or POSIX ACLs.

* **Administration**: how much administration is required for the system to
  function?
    - 0: Frequent.
    - 1: Infrequent.
    - 2: Essentially none.

* **Simplicity**: how understandable is the system as a whole?
    - 0: Very complex.
    - 1: Understandable with some study.
    - 2: Very simple, easy to predict.

* **Visibility**: how much visibility is available into processes within the
  system?
    - 0: Total black box.
    - 1: Basic logging.
    - 2: CLI tooling.
    - 3: Exportable metrics (e.g. prometheus).

## Evaluations

With the rubric defined, let's start actually working through our options! There
are many, many different possibilities, so this may not be an exhaustive list.

### [Ceph](https://docs.ceph.com/en/latest/cephfs/index.html)

> The Ceph File System, or CephFS, is a POSIX-compliant file system built on
> top of Ceph’s distributed object store, RADOS. CephFS endeavors to provide a
> state-of-the-art, multi-use, highly available, and performant file store for
> a variety of applications, including traditional use-cases like shared home
> directories, HPC scratch space, and distributed workflow shared storage. 

- Hackability: 2. Very active community, but it's C++.
- Documentation: 2. Hella docs, very daunting.
- Transience: 0. Adding hosts seems to require multiple configuration steps.
- Priority: 1. There is fine-tuning on a per-host basis.
- Caching: 1. Clients can cache both metadata and block data.
- Conflicts: 1. The FS behaves as much like a real FS as possible.
- Consistency: 0. System is CP.
- POSIX: 2. Fully POSIX compliant.
- Scale: 2. Cluster can grow without any real bounds.
- Failure: 3. There's no indication anywhere that Ceph goes into any kind of cluster-wide failure mode.
- Limitations: 2. There are performance considerations with large directories, but no hard limits.
- Encryption: 0. None to speak of.
- Permissions: 2. POSIX ACLs supported.
- Administration: 1. This is a guess, but Ceph seems to be self-healing in general, but still needs hand-holding in certain situations (adding/removing nodes, etc...)
- Simplicity: 0. There are many moving pieces, as well as many different kinds of processes and entities.
- Visibility: 3. Lots of tooling to dig into the state of the cluster, as well as a prometheus module.

TOTAL: 22

#### Comments

Ceph has been recommended to me by a few people. It is clearly a very mature
project, though that maturity has brought with it a lot of complexity. A lot of
the complexity of Ceph seems to be rooted in its strong consistency guarantees,
which I'm confident it fulfills well, but are not really needed for the
use-case I'm interested in. I'd prefer a simpler, eventually consistent,
system. It's also not clear to me that Ceph would even perform very well in my
use-case as it seems to want an actual datacenter deployment, with beefy
hardware and hosts which are generally close together.

### [GlusterFS](https://docs.gluster.org/en/latest/)

> GlusterFS is a scalable network filesystem suitable for data-intensive tasks
> such as cloud storage and media streaming. GlusterFS is free and open source
> software and can utilize common off-the-shelf hardware. 

- Hackability: 2. Mostly C code, but there is an active community.
- Documentation: 2. Good docs.
- Transience: 0. New nodes cannot add themselves to the pool.
- Priority: 0. Data is distributed based on consistent hashing algo, nothing else.
- Caching: 1. Docs mention client-side caching layer.
- Conflicts: 0. File becomes frozen, manual intervention is needed to save it.
- Consistency: 0. Gluster aims to be fully consistent.
- POSIX: 2. Fully POSIX compliant.
- Scale: 2. No apparent limits.
- Failure: 3. Clients determine on their own whether or not they have a quorum for a particular sub-volume.
- Limitations: 2. Limited by the file system underlying each volume, I think.
- Encryption: 2. Encryption can be done on the volume level, each user could have a private volume.
- Permissions: 2. ACL checking is enforced on the server-side, but requires syncing of users and group membership across servers.
- Administration: 1. Beyond adding/removing nodes the system is fairly self-healing.
- Simplicity: 1. There's only one kind of server process, and the configuration of volumes is is well documented and straightforward.
- Visibility: 3. Prometheus exporter available.

TOTAL: 23

#### Comments

GlusterFS was my initial choice when I did a brief survey of DFSs for this
use-case. However, after further digging into it I think it will suffer the
same ultimate problem as CephFS: too much consistency for a wide-area
application like I'm envisioning. The need for syncing user/groups across
machines as actual system users is also cumbersome enough to make it not a
great choice.

### [MooseFS](https://moosefs.com/)

> MooseFS is a Petabyte Open Source Network Distributed File System. It is easy
> to deploy and maintain, highly reliable, fault tolerant, highly performing,
> easily scalable and POSIX compliant.
>
> MooseFS spreads data over a number of commodity servers, which are visible to
> the user as one resource. For standard file operations MooseFS acts like
> ordinary Unix-like file system.

- Hackability: 2. All C code, pretty dense, but backed by a company.
- Documentation: 2. There's a giant PDF you can read through like a book. I
  guess that's.... good?
- Transience: 0. Nodes must be added manually.
- Priority: 1. There's "Storage Classes".
- Caching: 1. Caching is done on the client, and there's some synchronization
  with the master server around it.
- Conflicts: 1. Both update operations will go through.
- Consistency: 0. Afaict it's a fully consistent system, with a master server
  being used to synchronize changes.
- POSIX: 2. Fully POSIX compliant.
- Scale: 2. Cluster can grow without any real bounds.
- Failure: 1. If the master server is unreachable then the client can't
  function.
- Limitations: 2. Limits are very large, effectively no limit.
- Encryption: 0. Docs make no mention of encryption.
- Permissions: 1. Afaict permissions are done by the OS on the fuse mount.
- Administration: 1. It seems that if the topology is stable there shouldn't be
  much going on.
- Simplicity: 0. There are many moving pieces, as well as many different kinds of processes and entities.
- Visibility: 2. Lots of cli tooling, no prometheus metrics that I could find.

TOTAL: 17

Overall MooseFS seems to me like a poor-developer's Ceph. It can do exactly the
same things, but with less of a community around it. The sale's pitch and
feature-gating also don't ingratiate it to me. The most damning "feature" is the
master metadata server, which acts as a SPOF and only sort of supports
replication (but not failover, unless you get Pro).

## Cutting Room Floor

The following projects were intended to be reviewed, but didn't make the cut for
various reasons.

* Tahoe-LAFS: The FUSE mount (which is actually an SFTP mount) doesn't support
  mutable files.

* HekaFS: Doesn't appear to exist anymore(?)

* IPFS-cluster: Doesn't support sharding.

* MinFS: Seems to only work off S3, no longer maintained anyway.

* DRDB: Linux specific, no mac support.

* BeeGFS: No mac support (I don't think? I couldn't find any indication it
  supports macs at any rate).

* NFS: No support for sharding the dataset.

## Conclusions

Going through the featuresets of all these different projects really helped me
focus in on how I actually expect this system to function, and a few things
stood out to me:

* Perfect consistency is not a goal, and is ultimately harmful for this
  use-case. The FS needs to propagate changes relatively quickly, but if two
  different hosts are updating the same file it's not necessary to synchronize
  those updates like a local filesystem would; just let one changeset clobber
  the other and let the outer application deal with coordination.

* Permissions are extremely important, and yet for all these projects are
  generally an afterthought. In a distributed setting we can't rely on the OS
  user/groups of a host to permission read/write access. Instead that must be
  done primarily via e2e encryption.

* Transience is not something most of these project expect, but is a hard
  requirement of this use-case. In the long run we need something which can be
  run on home hardware on home ISPs, which is not reliable at all. Hosts need to
  be able to flit in and out of existence, and the cluster as a whole needs to
  self-heal through that process.

In the end, it may be necessary to roll our own project for this, as I don't
think any of the existing distributed file systems are suitable for what's
needed.

[mobile_nebula]: https://github.com/cryptic-io/mobile_nebula
[nix]: https://nixos.org/manual/nix/stable/
