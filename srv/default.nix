{
  bash,
  go,
  buildGoModule,
  writeScript,
  writeText,
  stdenv,

  config,
  staticBuild,
}: rec {

    init = writeText "mediocre-blog-srv-init" ''

      export MEDIOCRE_BLOG_DATA_DIR="${config.dataDir}"
      mkdir -p "${config.dataDir}"

      # mailing list
      export MEDIOCRE_BLOG_ML_SMTP_ADDR="${config.mlSMTPAddr}"
      export MEDIOCRE_BLOG_ML_SMTP_AUTH="${config.mlSMTPAuth}"
      export MEDIOCRE_BLOG_ML_PUBLIC_URL="${config.mlPublicURL}"

      # redis
      export MEDIOCRE_BLOG_REDIS_PROTO=unix
      export MEDIOCRE_BLOG_REDIS_ADDR="${config.redisListenPath}"

      # pow
      export MEDIOCRE_BLOG_POW_SECRET="${config.powSecret}"

      # static proxy
      if [ "${config.staticProxyURL}" == "" ]; then
        export MEDIOCRE_BLOG_STATIC_DIR="${staticBuild}"
      else
        export MEDIOCRE_BLOG_STATIC_URL="${config.staticProxyURL}"
      fi

      # listening
      export MEDIOCRE_BLOG_LISTEN_PROTO="${config.listenProto}"
      export MEDIOCRE_BLOG_LISTEN_ADDR="${config.listenAddr}"
    '';

    build = buildGoModule {
        pname = "mediocre-blog-srv";
        version = "dev";
        src = ./src;
        vendorSha256 = "sha256-/F62WVkI50woo5J0xZOAn0g0WWkDna4wIBeVvbhAGzs=";

        # disable tests
        checkPhase = '''';
    };

    bin = writeScript "mediocre-blog-srv-bin" ''
        #!${bash}
        source ${init}
        exec ${build}/bin/mediocre-blog
    '';

    shell = stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell";
        buildInputs = [ go build ];
        shellHook = ''
          source ${init}
          cd src
        '';
    };
}
