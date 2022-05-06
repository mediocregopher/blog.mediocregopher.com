{
  buildGoModule,
  writeScript,
  writeScriptBin,
  stdenv,

  config,
  staticBuild,
}: rec {

    env = ''

      export MEDIOCRE_BLOG_DATA_DIR=${config.dataDir}

      # mailing list
      export MEDIOCRE_BLOG_ML_SMTP_ADDR=${config.mlSMTPAddr}
      export MEDIOCRE_BLOG_ML_SMTP_AUTH='${config.mlSMTPAuth}'
      export MEDIOCRE_BLOG_ML_PUBLIC_URL=${config.mlPublicURL}

      # redis
      export MEDIOCRE_BLOG_REDIS_PROTO=unix
      export MEDIOCRE_BLOG_REDIS_ADDR=${config.redisListenPath}

      # pow
      export MEDIOCRE_BLOG_POW_SECRET=${config.powSecret}

      # static proxy
      if [ "${config.staticProxyURL}" == "" ]; then
        export MEDIOCRE_BLOG_STATIC_DIR=${staticBuild}
      else
        export MEDIOCRE_BLOG_STATIC_URL=${config.staticProxyURL}
      fi

      # listening
      export MEDIOCRE_BLOG_LISTEN_PROTO=${config.listenProto}
      export MEDIOCRE_BLOG_LISTEN_ADDR=${config.listenAddr}
    '';

    build = buildGoModule {
        pname = "mediocre-blog-srv";
        version = "dev";
        src = ./.;
        vendorSha256 = "02szg1lisfjk8pk9pflbyv97ykg9362r4fhd0w0p2a7c81kf9b8y";

        # disable tests
        checkPhase = '''';
    };

    bin = writeScript "mediocre-blog-srv-bin" ''
        #!/bin/sh
        mkdir -p "${config.dataDir}"
        source ${env}
        exec ${build}/bin/mediocre-blog
    '';

    shell = stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell";
        shellHook = ''
          source ${env}
        '';
    };
}
