---
title: >-
    Building AppImages with Nix
description: >-
    With some process trees thrown in there for fun.
series: nebula
tags: tech
---

It's been a bit since I've written an update on the cryptic nebula project,
almost 5 months (since [this post][lastnix], which wasn't officially part of the
blog series but whatever). Since then it's switched names to "cryptic-net", and
that we would likely use [MinIO](https://min.io/) as our network storage
service, but neither of those is the most interesting update.

The project had been stalled because of a lack of a build system which could
fulfill the following requirements:

* Network configuration (static IP, VPN certificates) of individual hosts is
  baked into the binary they run.

* Binaries are completely static; no external dependencies need to exist on the
  host in order to run them.

* Each binary runs a composition of multiple sub-services, each being a separate
  sub-process, and all of them having been configured to work together (with
  some possible glue code on our side) to provide the features we want.

* The builder itself should be deterministic; no matter where it runs it should
  produce the same binary given the same input parameters.

Lacking such a build system we're not able to distribute cryptic-net in a way
which "just works"; it would require some kind of configuration, or some kind of
runtime environment to be set up, both of which would be a pain for users. And
lacking a definite build system makes it difficult to move forward on any other
aspect of a project, as it's not clear what may need to be redone in the future
when the build system is decided upon.

## Why not nix-bundle?

My usage of [nix-bundle][nix-bundle] in a [previous post][lastnix] was an
attempt at fulfilling these requirements. Nix in general does very well in
fulfilling all but the second requirement, and nix-bundle was supposed to
fulfill even that by packaging a nix derivation into a static binary.

And all of this it did! Except that the mechanism of nix-bundle is a bit odd.
The process of a nix-bundle'd binary jails itself within a chroot, which it then
uses to fake the `/nix/store` path which nix built binaries expect to exist.

This might work in a lot of cases, but it did not work in ours. For one, [nebula
can't create its network interface when run from inside
nix-bundle's chroot][nix-bundle-issue]. For another, being run in a chroot means
there's going to be strange restrictions on what our binary is able to do and
not.

## AppImage

What we really needed was an [AppImage][appimage]. AppImages are static binaries
which can bundle complex applications, even those which don't expect to be
bundled into single binaries. In this way the end result is the same as
nix-bundle, but the mechanism AppImage uses is different and places far fewer
restrictions on what we can and can't do with our program.

## Building Sub-Services Statically with Nix

It's probably possible to use nix to generate an AppImage which has the
`/nix/store` built into it, similar to what nix-bundle does, and therefore not
worry about whether the binaries it's bundling are static or not. But if your
services are written in sane languages it's not that difficult to build them
statically and dodge the issue.

For example, here is how you build a go binary statically:

```
{
    buildGoModule,
    fetchFromGitHub,
}:
    buildGoModule rec {
        pname = "nebula";
        version = "1.4.0";

        src = fetchFromGitHub {
            owner = "slackhq";
            repo = pname;
            rev = "v${version}";
            sha256 = "lu2/rSB9cFD7VUiK+niuqCX9CI2x+k4Pi+U5yksETSU=";
        };

        vendorSha256 = "p1inJ9+NAb2d81cn+y+ofhxFz9ObUiLgj+9cACa6Jqg=";

        doCheck = false;

        subPackages = [ "cmd/nebula" "cmd/nebula-cert" ];

        CGO_ENABLED=0;
        tags = [ "netgo" ];
        ldflags = [
            "-X main.Build=${version}"
            "-w"
            "-extldflags=-static"
        ];
    };
```

And here's how to statically build a C binary:

```
{
    stdenv,
    glibcStatic, # e.g. pkgs.glibc.static
}:
    stdenv.mkDerivation rec {
        pname = "dnsmasq";
        version = "2.85";

        src = builtins.fetchurl {
          url = "https://www.thekelleys.org.uk/dnsmasq/${pname}-${version}.tar.xz";
          sha256 = "sha256-rZjTgD32h+W5OAgPPSXGKP5ByHh1LQP7xhmXh/7jEvo=";
        };

        nativeBuildInputs = [ glibcStatic ];

        makeFlags = [
            "LDFLAGS=-static"
            "DESTDIR="
            "BINDIR=$(out)/bin"
            "MANDIR=$(out)/man"
            "LOCALEDIR=$(out)/share/locale"
        ];
    };
```

The derivations created by either of these expressions can be plugged right into
the `pkgs.buildEnv` used to create the AppDir (see AppDir section below).

## Process Manager

An important piece of the puzzle for getting cryptic-net into an AppImage was a
process manager. We need something which can run multiple service processes
simultaneously, restart processes which exit unexpectedly, gracefully handle
shutting down all those processes, and coalesce the logs of all processes into a
single stream.

There are quite a few process managers out there which could fit the bill, but
finding any which could be statically compiled ended up not being an easy task.
In the end I decided to see how long it would take me to implement such a
program in go, and hope it would be less time than it would take to get
`circus`, a python program, bundled into the AppImage.

2 hours later, [pmux][pmux] was born! Check it out. It's a go program so
building it looks pretty similar to the nebula builder above, so I won't repeat
it. However I will show the configuration we're using for it within the
AppImage, to show how it ties all the processes together:

```yaml
processes:
    - name: nebula
      cmd: bin/nebula
      args:
        - "-config"
        - etc/nebula/nebula.yml

    - name: dnsmasq
      cmd: bin/dnsmasq
      args:
        - "-d"
        - "-C"
        - ${dnsmasq}/etc/dnsmasq/dnsmasq.conf
```

## AppDir -> AppImage

Generating an AppImage requires an AppDir. An AppDir is a directory which
contains all files required by a program, rooted to the AppDir. For example, if
the program expects a file to be at `/etc/some/conf`, then that file should be
places in the AppDir at `<AppDir-path>/etc/some/conf`.

[These docs](https://docs.appimage.org/packaging-guide/manual.html#ref-manual)
were very helpful for me in figuring out how to construct the AppDir. I then
used the `pkgs.buildEnv` utility to create an AppDir derivation containing
everything cryptic-net needs to run:

```
    appDir = pkgs.buildEnv {
        name = "cryptic-net-AppDir";
        paths = [

            # real directory containing non-built files, e.g. the pmux config
            ./AppDir

            # static binary derivations shown previously
            nebula
            dnsmasq
            pmux
        ];
    };
```

Once the AppDir is built one needs to use `appimagetool` to turn it into an
AppImage. There is an `appimagetool` build in the standard nixpkgs, but
unfortunately it doesn't seem to actually work...

Luckily nix-bundle is working on AppImage support, and includes a custom build
of `appimagetool` which does work!

```
{
    fetchFromGitHub,
    callPackage,
}: let
    src = fetchFromGitHub {
        owner = "matthewbauer";
        repo = "nix-bundle";
        rev = "223f4ffc4179aa318c34dc873a08cb00090db829";
        sha256 = "0pqpx9vnjk9h24h9qlv4la76lh5ykljch6g487b26r1r2s9zg7kh";
    };
in
    callPackage "${src}/appimagetool.nix" {}
```

Using `callPackage` on this expression will give you a functional `appimagetool`
derivation. From there's it's a simple matter of writing a derivation which
generates the AppImage from a created AppDir:

```
{
    appDir,
    appimagetool,
}:
    pkgs.stdenv.mkDerivation {
        name = "cryptic-net-AppImage";

        src = appDir;
        buildInputs = [ appimagetool ];
        ARCH = "x86_64"; # required by appimagetool

        builder = builtins.toFile "build.sh" ''
            source $stdenv/setup
            cp -rL "$src" buildAppDir
            chmod +w buildAppDir -R
            mkdir $out

            appimagetool cryptic-net "$out/cryptic-net-bin"
        '';
    }
```

Running that derivation deterministically spits out a binary at
`result/cryptic-net-bin` which can be executed and run immediately, on any
system using the `x86_46` CPU architecture.

## Fin

I'm extremely hyped to now have the ability to generate binaries for cryptic-net
that people can _just run_, without them worrying about which sub-services that
binary is running under-the-hood. From a usability perspective it's way nicer
than having to tell people to "install docker" or "install nix", and from a dev
perspective we have a really solid foundation on which to build a quite complex
application.

[lastnix]: {% post_url 2021-04-22-composing-processes-into-a-static-binary-with-nix %}
[nix-bundle]: https://github.com/matthewbauer/nix-bundle
[nix-bundle-issue]: https://github.com/matthewbauer/nix-bundle/issues/78
[appimage]: https://appimage.org/
[pmux]: https://github.com/cryptic-io/pmux
