---
title: >-
    Composing Processes Into a Static Binary With Nix
description: >-
    Goodbye, docker-compose!
tags: tech
---

It's pretty frequent that one wants to use a project that requires multiple
processes running. For example, a small web api which uses some database to
store data in, or a networking utility which has some monitoring process which
can be run alongside it.

In these cases it's extremely helpful to be able to compose these disparate
processes together into a single process. From the user's perspective it's much
nicer to only have to manage one process (even if it has hidden child
processes). From a dev's perspective the alternatives are: finding libraries in
the same language which do the disparate tasks and composing them into the same
process via import, or (if such libraries don't exist, which is likely)
rewriting the functionality of all processes into a new, monolithic project
which does everything; a huge waste of effort!

## docker-compose

A tool I've used before for process composition is
[docker-compose][docker-compose]. While it works well for composition, it
suffers from the same issues docker in general suffers from: annoying networking
quirks, a questionable security model, and the need to run the docker daemon.
While these issues are generally surmountable for a developer or sysadmin, they
are not suitable for a general-purpose project which will be shipped to average
users.

## nix-bundle

Enter [nix-bundle][nix-bundle]. This tools will take any [nix][nix] derivation
and construct a single static binary out of it, a la [AppImage][appimage].
Combined with a process management tool like [circus][circus], nix-bundle
becomes a very useful tool for composing processes together!

To demonstrate this, we'll be looking at putting together a project I wrote
called [markov][markov], a simple REST API for building [markov
chains][markov-chain] which is written in [go][golang] and backed by
[redis][redis].

## Step 1: Building Individual Components

Step one is to get [markov][markov] and its dependencies into a state where it
can be run with [nix][nix]. Doing this is fairly simple, we merely use the
`buildGoModule` function:

```
pkgs.buildGoModule {
    pname = "markov";
    version = "618b666484566de71f2d59114d011ff4621cf375";
    src = pkgs.fetchFromGitHub {
        owner = "mediocregopher";
        repo = "markov";
        rev = "618b666484566de71f2d59114d011ff4621cf375";
        sha256 = "1sx9dr1q3vr3q8nyx3965x6259iyl85591vx815g1xacygv4i4fg";
    };
    vendorSha256 = "048wygrmv26fsnypsp6vxf89z3j0gs9f1w4i63khx7h134yxhbc6";
}
```

This expression results in a derivation which places the markov binary at
`bin/markov`.

The other component we need to run markov is [redis][redis], which conveniently
is already packaged in nixpkgs as `pkg.redis`.

## Step 2: Composing Using Circus

[Circus][circus] can be configured to run multiple processes at the same time.
It will collect the stdout/stderr logs of these processes and combine them into
a single stream, or write them to log files. If any processes fail circus will
automatically restart them. It has a simple configuration and is, overall, a
great tool for a simple project like this.

Circus also comes pre-packed in nixpkgs, so we don't need to do anything to
actually build it. We only need to configure it. To do this we'll write a bash
script which generates the configuration on-the-fly, and then runs the process
with that configuration.

This script is going to act as the "frontend" for our eventual static binary;
the user will pass in configuration parameters to this script, and this script
will translate those into the appropriate configuration for all sub-process
(markov, redis, circus). For this demo we won't go nuts with the configuration,
we'll just expose the following:

* `MARKOV_LISTEN_ADDR`: Address REST API will listen on (defaults to
  `localhost:8000`).

* `MARKOV_TIMEOUT`: Expiration time of each link of the chain (defaults to 720
  hours).

* `MARKOV_DATA_DIR`: Directory where data will be stored (defaults to current
  working directory).

The bash script will take these params in as environment variables. The nix
expression to generate the bash script, which we'll call our entrypoint script,
will look like this (assumes that the expression to generate `bin/markov`,
defined above, is set to the `markov` variable):

```
pkgs.writeScriptBin "markov" ''
    #!${pkgs.stdenv.shell}

    # On every run we create new, temporary, configuration files for redis and
    # circus. To do this we create a new config directory.
    markovCfgDir=$(${pkgs.coreutils}/bin/mktemp -d)
    echo "generating configuration to $markovCfgDir"

    cat >$markovCfgDir/redis.conf <<EOF
    save ""
    dir "''${MARKOV_DATA_DIR:-$(pwd)}"
    appendonly yes
    appendfilename "markov.data"
    EOF

    cat >$markovCfgDir/circus.ini <<EOF

    [circus]

    [watcher:markov]
    cmd = ${markov}/bin/markov \
        -listenAddr ''${MARKOV_LISTEN_ADDR:-localhost:8000} \
        -timeout ''${MARKOV_TIMEOUT:-720}
    numprocesses = 1

    [watcher:redis]
    cmd = ${pkgs.redis}/bin/redis-server $markovCfgDir/redis.conf
    numprocesses = 1
    EOF

    exec ${pkgs.circus}/bin/circusd $markovCfgDir/circus.ini
'';
```

By `nix-build`ing this expression we end up with a derivation with
`bin/markov`, and running that should result in the following output:

```
generating configuration to markov.VLMPwqY
2021-04-22 09:27:56 circus[181906] [INFO] Starting master on pid 181906
2021-04-22 09:27:56 circus[181906] [INFO] Arbiter now waiting for commands
2021-04-22 09:27:56 circus[181906] [INFO] markov started
2021-04-22 09:27:56 circus[181906] [INFO] redis started
181923:C 22 Apr 2021 09:27:56.063 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
181923:C 22 Apr 2021 09:27:56.063 # Redis version=6.0.6, bits=64, commit=00000000, modified=0, pid=181923, just started
181923:C 22 Apr 2021 09:27:56.063 # Configuration loaded
...
```

The `markov` server process doesn't have many logs, unfortunately, but redis'
logs at least work well, and doing a `curl localhost:8000` results in the
response from the `markov` server.

At this point our processes are composed using circus, let's now bundle it all
into a single static binary!

## Step 3: nix-bundle

The next step is to run [nix-bundle][nix-bundle] on the entrypoint expression,
and nix-bundle will compile all dependencies (including markov, redis, and
circus) into a single archive file, and make that file executable. When the
archive is executed it will run our entrypoint script directly.

Getting nix-bundle is very easy, just use nix-shell!

```
nix-shell -p nix-bundle
```

This will open a shell where the `nix-bundle` binary is available on your path.
From there just run the following to construct the binary (this assumes that the
nix code described so far is stored in `markov.nix`, the full source of which
will be linked to at the end of this post):

```
nix-bundle '((import ./markov.nix) {}).entrypoint' '/bin/markov'
```

The resulting binary is called `markov`, and is 89MB. The size is a bit jarring,
considering the simplicity of the functionality, but it could probably be
trimmed by using a different process manager than circus (which requires
bundling an entire python runtime into the binary).

Running the binary directly as `./markov` produces the same result as when we
ran the entrypoint script earlier. Success! We have bundled multiple existing
processes into a single, opaque, static binary. Installation of this binary is
now as easy as copying it to any linux machine and running it.

## Bonus Step: nix'ing nix-bundle

Installing and running [nix-bundle][nix-bundle] manually is _fine_, but it'd be even better if
that was defined as part of our nix setup as well. That way any new person
wouldn't have to worry about that step, and still get the same deterministic
output from the build.

Unfortunately, we can't actually run `nix-bundle` from within a nix build
derivation, as it requires access to the nix store and that can't be done (or at
least I'm not on that level yet). So instead we'll have to settle for defining
the `nix-bundle` binary in nix and then using a `Makefile` to call it.

Defining a `nix-bundle` expression is easy enough:

```
    nixBundleSrc = pkgs.fetchFromGitHub {
        owner = "matthewbauer";
        repo = "nix-bundle";
        rev = "8e396533ef8f3e8a769037476824d668409b4a74";
        sha256 = "1lrq0990p07av42xz203w64abv2rz9xd8jrzxyvzzwj7vjj7qwyw";
    };

    nixBundle = (import "${nixBundleSrc}/release.nix") {
        nixpkgs' = pkgs;
    };
```

Then the Makefile:

```make
bundle:
	nix-build markov.nix -A nixBundle
	./result/bin/nix-bundle '((import ./markov.nix) {}).entrypoint' '/bin/markov'
```

Now all a developer needs to rebuild the project is to do `make` within the
directory, while also having nix set up. The result will be a deterministically
built, static binary, encompassing multiple processes which will all work
together behind the scenes. This static binary can be copied to any linux
machine and run there without any further installation steps.

How neat is that!

The final source files used for this project can be found below:

* [markov.nix](/assets/markov/markov.nix.html)
* [Makefile](/assets/markov/Makefile.html)

[nix]: https://nixos.org/manual/nix/stable/
[nix-bundle]: https://github.com/matthewbauer/nix-bundle
[docker-compose]: https://docs.docker.com/compose/
[appimage]: https://appimage.org/
[circus]: https://circus.readthedocs.io/en/latest/
[markov]: https://github.com/mediocregopher/markov
[markov-chain]: https://en.wikipedia.org/wiki/Markov_chain
[golang]: https://golang.org/
[redis]: https://redis.io/
