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
    ] ++ (
        if config.staticProxyURL == ""
        then [ "-static-dir=${staticBuild}" ]
        else [ "-static-proxy-url=${config.staticProxyURL}" ]
    );

    build = pkgs.buildGoModule {
        pname = "mediocre-blog-srv";
        version = "dev";
        src = ./.;
        vendorSha256 = "08wv94yv2wmlxzmanw551gixc8v8nl6zq2m721ig9nl3r540x46f";
    };

    bin = pkgs.writeScript "mediocre-blog-srv-bin" ''
        #!/bin/sh
        exec ${build}/bin/mediocre-blog ${toString opts}
    '';

    runScript = pkgs.writeScriptBin "run-mediocre-blog" ''
        go run ./cmd/mediocre-blog/main.go ${toString opts}
    '';

    runMailingListCLIScript = pkgs.writeScriptBin "run-mailinglist-cli" ''
        go run ./cmd/mailinglist-cli/main.go ${toString mailingListOpts} $@
    '';

    shell = pkgs.stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell";
        buildInputs = [ pkgs.go runScript runMailingListCLIScript ];
    };

}
