{
  bash,
  go,
  buildGoModule,
  writeScript,
  writeText,
  stdenv,

  config,
}: rec {

    init = writeText "mediocre-blog-srv-init" ''

      export MEDIOCRE_BLOG_DATA_DIR="${config.dataDir}"

      # mailing list
      export MEDIOCRE_BLOG_ML_SMTP_ADDR="${config.mlSMTPAddr}"
      export MEDIOCRE_BLOG_ML_SMTP_AUTH="${config.mlSMTPAuth}"
      export MEDIOCRE_BLOG_ML_PUBLIC_URL="${config.mlPublicURL}"

      # redis
      export MEDIOCRE_BLOG_REDIS_PROTO=unix
      export MEDIOCRE_BLOG_REDIS_ADDR="${config.redisListenPath}"

      # pow
      export MEDIOCRE_BLOG_POW_SECRET="${config.powSecret}"

      # http
      export MEDIOCRE_BLOG_LISTEN_PROTO="${config.httpListenProto}"
      export MEDIOCRE_BLOG_LISTEN_ADDR="${config.httpListenAddr}"
      export MEDIOCRE_BLOG_HTTP_AUTH_USERS='${builtins.toJSON config.httpAuthUsers}'
      export MEDIOCRE_BLOG_HTTP_AUTH_RATELIMIT='${config.httpAuthRatelimit}'
    '';

    build = buildGoModule {
        pname = "mediocre-blog-srv";
        version = "dev";
        src = ./src;
        vendorSha256 = "1s5jhis1a2y7m50k29ap7kd0h4bgc3dzy1f9dqf5jrz8n27f3i87";

        # disable tests
        checkPhase = '''';
    };

    bin = writeScript "mediocre-blog-srv-bin" ''
        #!${bash}/bin/bash
        source ${init}
        exec ${build}/bin/mediocre-blog
    '';

    shell = stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell";
        buildInputs = [ go ];
        shellHook = ''source ${init}'';
    };

    shellWithBuild = stdenv.mkDerivation {
        name = "mediocre-blog-srv-shell-with-build";
        buildInputs = [ go build ];
        shellHook = ''source ${init}'';
    };
}
