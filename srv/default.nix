{pkgs, config, staticBuild}: rec {

    mailingListOpts = [
        "-ml-smtp-addr=${config.mlSMTPAddr}"
        "-ml-smtp-auth='${config.mlSMTPAuth}'"
        "-data-dir=${config.dataDir}"
        "-public-url=${config.publicURL}"
    ];

    opts = mailingListOpts ++ [
        "-pow-secret=${config.powSecret}"
        "-listen-proto=${config.listenProto}"
        "-listen-addr=${config.listenAddr}"
        "-redis-proto=unix"
        "-redis-addr=${config.redisListenPath}"
    ] ++ (
        if config.staticProxyURL == ""
        then [ "-static-dir=${staticBuild}" ]
        else [ "-static-proxy-url=${config.staticProxyURL}" ]
    );

    build = pkgs.buildGoModule {
        pname = "mediocre-blog-srv";
        version = "dev";
        src = ./.;
        vendorSha256 = "0c6j989q6r2q967gx90cl4l8skflkx2npmxd3f5l16bwj2ldw11j";

        # disable tests
        checkPhase = '''';
    };

    bin = pkgs.writeScript "mediocre-blog-srv-bin" ''
        #!/bin/sh
        mkdir -p "${config.dataDir}"
        exec ${build}/bin/mediocre-blog ${toString opts}
    '';

    runScript = pkgs.writeScriptBin "run-mediocre-blog" ''
        mkdir -p "${config.dataDir}"
        go run ./cmd/mediocre-blog/main.go ${toString opts}
    '';

    runMailingListCLIScript = pkgs.writeScriptBin "run-mailinglist-cli" ''
        go run ./cmd/mailinglist-cli/main.go ${toString mailingListOpts} "$@"
    '';

    shell = pkgs.stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell";
        buildInputs = [ pkgs.go runScript runMailingListCLIScript ];
    };

}
