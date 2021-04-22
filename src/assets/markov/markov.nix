{
    pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/d50923ab2d308a1ddb21594ba6ae064cab65d8ae.tar.gz") {}
}:

rec {

    markov = pkgs.buildGoModule {
        pname = "markov";
        version = "618b666484566de71f2d59114d011ff4621cf375";
        src = pkgs.fetchFromGitHub {
            owner = "mediocregopher";
            repo = "markov";
            rev = "618b666484566de71f2d59114d011ff4621cf375";
            sha256 = "1sx9dr1q3vr3q8nyx3965x6259iyl85591vx815g1xacygv4i4fg";
        };
        vendorSha256 = "048wygrmv26fsnypsp6vxf89z3j0gs9f1w4i63khx7h134yxhbc6";
    };

    entrypoint = pkgs.writeScriptBin "markov" ''
        #!${pkgs.stdenv.shell}

        # On every run we create new, temporary, configuration files for redis and
        # circus. To do this we create a new config directory.
        markovCfgDir=$(${pkgs.coreutils}/bin/mktemp -d)
        echo "generating configuration to $markovCfgDir"

        ${pkgs.coreutils}/bin/cat >$markovCfgDir/redis.conf <<EOF
        save ""
        dir "''${MARKOV_DATA_DIR:-$(pwd)}"
        appendonly yes
        appendfilename "markov.data"
        EOF

        ${pkgs.coreutils}/bin/cat >$markovCfgDir/circus.ini <<EOF

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

    nixBundleSrc = pkgs.fetchFromGitHub {
        owner = "matthewbauer";
        repo = "nix-bundle";
        rev = "8e396533ef8f3e8a769037476824d668409b4a74";
        sha256 = "1lrq0990p07av42xz203w64abv2rz9xd8jrzxyvzzwj7vjj7qwyw";
    };

    nixBundle = (import "${nixBundleSrc}/release.nix") {
        nixpkgs' = pkgs;
    };
}

