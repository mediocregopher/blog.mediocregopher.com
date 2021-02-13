---
title: >-
    Building gomobile Using Nix
description: >-
    Harder than I thought it would be!
series: nebula
tags: tech
---

When I last left off with the nebula project I wanted to [nix][nix]-ify the
build process for Cryptic's [mobile_nebula][mobile_nebula] fork. While I've made
progress on the overall build, one particular bit of it really held me up, so
I'm writing about that part here. I'll finish the full build at a later time.

## gomobile

[gomobile][gomobile] is a toolkit for the go programming language to allow for
running go code on Android and iOS devices. `mobile_nebula` uses `gomobile` to
build a simple wrapper around the nebula client that the mobile app can then
hook into.

This means that in order to nix-ify the entire `mobile_nebula` project I first
need to nix-ify `gomobile`, and since there isn't (at time of writing) an
existing package for `gomobile` in the nixpkgs repo, I had to roll my own.

I started with a simple `buildGoModule` nix expression:

```
pkgs.buildGoModule {
    pname = "gomobile";
    version = "unstable-2020-12-17";
    src = pkgs.fetchFromGitHub {
        owner = "golang";
        repo = "mobile";
        rev = "e6ae53a27f4fd7cfa2943f2ae47b96cba8eb01c9";
        sha256 = "03dzis3xkj0abcm4k95w2zd4l9ygn0rhkj56bzxbcpwa7idqhd62";
    };
    vendorSha256 = "1n1338vqkc1n8cy94501n7jn3qbr28q9d9zxnq2b4rxsqjfc9l94";
}
```

The basic idea here is that `buildGoModule` will acquire a specific revision of
the `gomobile` source code from github, then attempt to build it. However,
`gomobile` is a special beast in that it requires a number of C/C++ libraries in
order to be built. I discovered this upon running this expression, when I
received this error:

```
./work.h:12:10: fatal error: GLES3/gl3.h: No such file or directory
   12 | #include <GLES3/gl3.h> // install on Ubuntu with: sudo apt-get install libegl1-mesa-dev libgles2-mesa-dev libx11-dev
```

This stumped me for a bit, as I couldn't figure out a) the "right" place to
source the `GLES3` header file from, and b) how to properly hook that into the
`buildGoModule` expression. My initial attempts involved trying to include
versions of the header file from my `androidsdk` nix package which I had already
gotten (mostly) working, but the version which ships there appears to expect to
be using clang. `cgo` (go's compiler which is used for C/C++ interop) only
supports gcc, so that strategy failed.

I didn't like having to import the header file from `androidsdk` anyway, as it
meant that my `gomobile` would only work within the context of the
`mobile_nebula` project, rather than being a standalone utility.

## nix-index

At this point I flailed around some more trying to figure out where to get this
header file from. Eventually I stumbled on the [nix-index][nix-index] project,
which implements something similar to the `locate` utility on linux: you give it
a file pattern, and it searches your active nix channels for any packages which
provide a file matching that pattern.

Since nix is amazing it's not actually necessary to install `nix-index`, I
simply start up a shell with the package available using `nix-shell -p
nix-index`. On first run I needed to populate the index by running the
`nix-index` command, which took some time, but after that finding packages which
provide the file I need is as easy as:

```
> nix-shell -p nix-index
[nix-shell:/tmp]$ nix-locate GLES3/gl3.h
(zulip.out)                                      82,674 r /nix/store/wbfw7w2ixdp317wip77d4ji834v1k1b9-libglvnd-1.3.2-dev/include/GLES3/gl3.h
libglvnd.dev                                     82,674 r /nix/store/pghxzmnmxdcarg5bj3js9csz0h85g08m-libglvnd-1.3.2-dev/include/GLES3/gl3.h
emscripten.out                                   82,666 r /nix/store/x3c4y2h5rn1jawybk48r6glzs1jl029s-emscripten-2.0.1/share/emscripten/system/include/GLES3/gl3.h
```

So my mystery file is provided by a few packages, but `libglvnd.dev` stood out
to me as it's also the pacman package which provides the same file in my real
operating system:

```
> yay -Qo /usr/include/GLES3/gl3.h
/usr/include/GLES3/gl3.h is owned by libglvnd 1.3.2-1
```

This gave me some confidence that this was the right track.

## cgo

My next fight was with `cgo` itself. Go's build process provides a few different
entry points for C/C++ compiler/linker flags, including both environment
variables and command-line arguments. But I wasn't using `go build` directly,
instead I was working through nix's `buildGoModule` wrapper. This added a huge
layer of confusion as all of nixpkgs is pretty terribly documented, so you
really have to just divine behavior from the [source][buildGoModule-source]
(good luck).

After lots of debugging (hint: `NIX_DEBUG=1`) I determined that all which is
actually needed is to set the `CGO_CFLAGS` variable within the `buildGoModule`
arguments. This would translate to the `CGO_CFLAGS` environment variable being
set during all internal commands, and whatever `go build` commands get used
would pick up my compiler flags from that.

My new nix expression looked like this:

```
pkgs.buildGoModule {
    pname = "gomobile";
    version = "unstable-2020-12-17";
    src = pkgs.fetchFromGitHub {
        owner = "golang";
        repo = "mobile";
        rev = "e6ae53a27f4fd7cfa2943f2ae47b96cba8eb01c9";
        sha256 = "03dzis3xkj0abcm4k95w2zd4l9ygn0rhkj56bzxbcpwa7idqhd62";
    };
    vendorSha256 = "1n1338vqkc1n8cy94501n7jn3qbr28q9d9zxnq2b4rxsqjfc9l94";

    CGO_CFLAGS = [
        "-I ${pkgs.libglvnd.dev}/include"
    ];
}
```

Running this produced a new error. Progress! The new error was:

```
/nix/store/p792j5f44l3f0xi7ai5jllwnxqwnka88-binutils-2.31.1/bin/ld: cannot find -lGLESv2
collect2: error: ld returned 1 exit status
```

So pretty similar to the previous issue, but this time the linker wasn't finding
a library file rather than the compiler not finding a header file. Once again I
used `nix-index`'s `nix-locate` command to find that this library file is
provided by the `libglvnd` package (as opposed to `libglvnd.dev`, which provided
the header file).

Adding `libglvnd` to the `CGO_CFLAGS` did not work, as it turns out that flags
for the linker `cgo` uses get passed in via `CGO_LDFLAGS` (makes sense). After
adding this new variable I got yet another error; this time `X11/Xlib.h` was not
able to be found. I repeated the process of `nix-locate`/add to `CGO_*FLAGS` a
few more times until all dependencies were accounted for. The new nix expression
looked like this:

```
pkgs.buildGoModule {
    pname = "gomobile";
    version = "unstable-2020-12-17";
    src = pkgs.fetchFromGitHub {
        owner = "golang";
        repo = "mobile";
        rev = "e6ae53a27f4fd7cfa2943f2ae47b96cba8eb01c9";
        sha256 = "03dzis3xkj0abcm4k95w2zd4l9ygn0rhkj56bzxbcpwa7idqhd62";
    };
    vendorSha256 = "1n1338vqkc1n8cy94501n7jn3qbr28q9d9zxnq2b4rxsqjfc9l94";

    CGO_CFLAGS = [
        "-I ${pkgs.libglvnd.dev}/include"
        "-I ${pkgs.xlibs.libX11.dev}/include"
        "-I ${pkgs.xlibs.xorgproto}/include"
        "-I ${pkgs.openal}/include"
    ];

    CGO_LDFLAGS = [
        "-L ${pkgs.libglvnd}/lib"
        "-L ${pkgs.xlibs.libX11}/lib"
        "-L ${pkgs.openal}/lib"
    ];
}
```

## Tests

The `CGO_*FLAGS` variables took care of all compiler/linker errors, but there
was one issue left: `buildGoModule` apparently runs the project's tests after
the build phase. `gomobile`'s tests were actually mostly passing, but some
failed due to trying to copy files around, which nix was having none of. After
some more [buildGoModule source][buildGoModule-source] divination I found that
if I passed an empty `checkPhase` argument it would skip the check phase, and
therefore skip running these tests.

## Fin!

The final nix expression looks like so:

```
pkgs.buildGoModule {
    pname = "gomobile";
    version = "unstable-2020-12-17";
    src = pkgs.fetchFromGitHub {
        owner = "golang";
        repo = "mobile";
        rev = "e6ae53a27f4fd7cfa2943f2ae47b96cba8eb01c9";
        sha256 = "03dzis3xkj0abcm4k95w2zd4l9ygn0rhkj56bzxbcpwa7idqhd62";
    };
    vendorSha256 = "1n1338vqkc1n8cy94501n7jn3qbr28q9d9zxnq2b4rxsqjfc9l94";

    CGO_CFLAGS = [
        "-I ${pkgs.libglvnd.dev}/include"
        "-I ${pkgs.xlibs.libX11.dev}/include"
        "-I ${pkgs.xlibs.xorgproto}/include"
        "-I ${pkgs.openal}/include"
    ];

    CGO_LDFLAGS = [
        "-L ${pkgs.libglvnd}/lib"
        "-L ${pkgs.xlibs.libX11}/lib"
        "-L ${pkgs.openal}/lib"
    ];

    checkPhase = "";
}
```

Once I complete the nix-ification of `mobile_nebula` I'll submit a PR to the
nixpkgs upstream with this, so that others can have `gomobile` available as
well!

[nix]: https://nixos.org/manual/nix/stable/
[mobile_nebula]: https://github.com/cryptic-io/mobile_nebula
[gomobile]: https://github.com/golang/mobile
[nix-index]: https://github.com/bennofs/nix-index
[buildGoModule-source]: https://github.com/NixOS/nixpkgs/blob/26117ed4b78020252e49fe75f562378063471f71/pkgs/development/go-modules/generic/default.nix
