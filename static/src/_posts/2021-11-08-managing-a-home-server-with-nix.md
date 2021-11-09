---
title: >-
    Managing a Home Server With Nix
description: >-
    Docker is for boomers.
tags: tech
---

My home server has a lot running on it. Some of it I've written about in this
blog previously, some of it I haven't. It's hosting this blog itself, even!

With all of these services comes management overhead, both in terms of managing
packages and configuration. I'm pretty strict about tracking packages and
configuration in version control, and backing up all state I care about in B2,
such that if, _at any moment_, the server is abducted by aliens, I won't have
lost much.

## Docker

Previously I accomplished this with docker. Each service ran in a container
under the docker daemon, with configuration files and state directories shared
in via volume shares. Configuration files could then be stored in a git repo,
and my `docker run` commands were documented in `Makefile`s, because that was
easy.

This approach had drawbacks, notably:

* Docker networking is a pain. To be fair I should have just used
  `--network=host` and dodged the issue, but I didn't.

* Docker images aren't actually deterministically built, so if I were to ever
  have to rebuild any of the images I was using it I couldn't be sure I'd end up
  with the same code as before. For some services this is actually a nagging
  security concern in the back of my head.

* File permissions with docker volumes are fucked.

* Who knows how long the current version of docker will support the old ass
  images and configuration system I'm using now. Probably not the next 10 years.
  And what if dockerhub goes away, or changes its pricing model?

* As previously noted, docker is for boomers.

## Nix

Nix is the new hotness, and it solves all of the above problems quite nicely.
I'm not going to get into too much detail about how nix works here (honestly I'm
not very good at explaining it), but suffice to say I'm switching everything
over, and this post is about how that actually looks in a practical sense.

For the most part I eschew things like [flakes][flakes],
[home-manager][home-manager], and any other frameworks built on nix. While the
framework of the day may come and go, the base nix language should remain
constant.

As before with docker, I have a single git repo being stored privately in a way
I'm confident is secure (which is necessary because it contains some secrets).

At the root of the repo there exists a `pkgs.nix` file, which looks like this:

```
{
  src ? builtins.fetchTarball {
    name = "nixpkgs-d50923ab2d308a1ddb21594ba6ae064cab65d8ae";
    url = "https://github.com/NixOS/nixpkgs/archive/d50923ab2d308a1ddb21594ba6ae064cab65d8ae.tar.gz";
    sha256 = "1k7xpymhzb4hilv6a1jp2lsxgc4yiqclh944m8sxyhriv9p2yhpv";
  },
}: (import src) {}
```

This file exists to provide a pinned version of `nixpkgs` which will get used
for all services. As long as I don't change this file the tools available to me
for building my services will remain constant forever, no matter what else
happens in the nix ecosystem.

Each directory in the repo corresponds to a service I run. I'll focus on a
particular service, [navidrome][navidrome], for now:

```bash
:: ls -1 navidrome
Makefile
default.nix
navidrome.toml
```

Not much to it!

### default.nix

The first file to look at is the `default.nix`, as that contains
all the logic. The overall file looks like this:

```
let

  pkgs = (import ../pkgs.nix) {};

in rec {

    entrypoint = ...;

    service = ...;

    install = ...;

}
```

The file describes an attribute set with three attributes, `entrypoint`,
`service`, and `install`. These form the basic pattern I use for all my
services; pretty much every service I manage has a `default.nix` which has
attributes corresponding to these.

#### Entrypoint

The first `entrypoint`, looks like this:

```
  entrypoint = pkgs.writeScript "mediocregopher-navidrome" ''
    #!${pkgs.bash}/bin/bash
    exec ${pkgs.navidrome}/bin/navidrome --configfile ${./navidrome.toml}
  '';
```

The goal here is to provide an executable which can be run directly, and which
will put together all necessary environment and configuration (`navidrome.toml`,
in this case) needed to run the service. Having the entrypoint split out into
its own target, as opposed to inlining it into the service file (defined next),
is convenient for testing; it allows you test _exactly_ what's going to happen
when running the service normally.

#### Service

`service` looks like this:

```
  service = pkgs.writeText "mediocregopher-navidrome-service" ''
    [Unit]
    Description=mediocregopher navidrome
    Requires=network.target
    After=network.target

    [Service]
    Restart=always
    RestartSec=1s
    User=mediocregopher
    Group=mediocregopher
    LimitNOFILE=10000

    # The important part!
    ExecStart=${entrypoint}

    # EXTRA DIRECTIVES ELIDED, SEE
    # https://www.navidrome.org/docs/installation/pre-built-binaries/

    [Install]
    WantedBy=multi-user.target
  '';
```

It's function is to produce a systemd service file. The service file will
reference the `entrypoint` which has already been defined, and in general does
nothing else.

#### Install

`install` looks like this:

```
  install = pkgs.writeScript "mediocregopher-navidrome-install" ''
    #!${pkgs.bash}/bin/bash
    sudo cp ${service} /etc/systemd/system/mediocregopher-navidrome.service
    sudo systemctl daemon-reload
    sudo systemctl enable mediocregopher-navidrome
    sudo systemctl restart mediocregopher-navidrome
  '';
```

This attribute produces a script which will install a systemd service on the
system it's run on. Assuming this is done in the context of a functional nix
environment and standard systemd installation it will "just work"; all relevant
binaries, configuration, etc, will all come along for the ride, and the service
will be running _exactly_ what's defined in my repo, everytime. Eat your heart
out, ansible!

Nix is usually used for building things, not _doing_ things, so it may seem
unusual for this to be here. But there's a very good reason for it, which I'll
get to soon.

### Makefile

While `default.nix` _could_ exist alone, and I _could_ just interact with it
directly using `nix-build` commands, I don't like to do that. Most of the reason
is that I don't want to have to _remember_ the `nix-build` commands I need. So
in each directory there's a `Makefile`, which acts as a kind of index of useful
commands. The one for navidrome looks like this:

```
install:
	$$(nix-build -A install --no-out-link)
```

Yup, that's it. It builds the `install` attribute, and runs the resulting script
inline. Easy peasy. Other services might have some other targets, like `init`,
which operate the same way but with different script targets.

## Nix Remotely

If you were waiting for me to explain _why_ the install target is in
`default.nix`, rather than just being in the `Makefile` (which would also make
sense), this is the part where I do that.

My home server isn't the only place where I host services, I also have a remote
host which runs some services. These services are defined in this same repo, in
essentially the same way as my local services. The only difference is in the
`Makefile`. Let's look at an example from my `maddy/Makefile`:

```
install-vultr:
	nix-build -A install --arg paramsFile ./vultr.nix
	nix-copy-closure -s ${VULTR} $$(readlink result)
	ssh -tt -q ${VULTR} $$(readlink result)
```

Vultr is the hosting company I'm renting the server from. Apparently I think I
will only ever have one host with them, because I just call it "vultr".

I'll go through this one line at a time. The first line is essentially the same
as the `install` line from my `navidrome` configuration, but with two small
differences: it takes in a parameters file containing the configuration
specific to the vultr host, and it's only _building_ the install script, not
running it.

The second line is the cool part. My remote host has a working nix environment
already, so I can just use `nix-copy-closure` to copy the `install` script to
it. Since the `install` script references the service file, which in turn
references the `entrypoint`, which in turn references the service binary itself,
and all of its configuration, _all_ of it will get synced to the remote host as
part of the `nix-copy-closure` command.

The third line runs the install script remotely. Since `nix-copy-closure`
already copied over all possible dependencies of the service, the end result is
a systemd service running _exactly_ as it would have if I were running it
locally.

All of this said, it's clear that provisioning this remote host in the first
place was pretty simple:

* Add my ssh key (done automatically by Vultr).
* Add my user to sudoers (done automatically by Vultr).
* Install single-user nix (two bash commands from
  [here](https://nixos.wiki/wiki/Nix_Installation_Guide#Stable_Nix)).

And that's literally it. No docker, no terraform, no kubernubernetes, no yaml
files... it all "just works". Will it ever require manual intervention? Yeah,
probably... I haven't defined uninstall or stop targets, for instance (though
that would be trivial to do). But overall, for a use-case like mine where I
don't need a lot, I'm quite happy.

That's pretty much the post. Hosting services at home isn't very difficult to
begin with, and with this pattern those of us who use nix can do so with greater
reliability and confidence going forward.

[flakes]: https://nixos.wiki/wiki/Flakes
[home-manager]: https://github.com/nix-community/home-manager
[navidrome]: https://github.com/navidrome/navidrome
