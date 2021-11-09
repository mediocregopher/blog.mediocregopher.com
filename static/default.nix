{pkgs}: rec {

    depInputs = [ pkgs.imagemagick pkgs.exiftool pkgs.bundler pkgs.bundix ];

    depShell = pkgs.stdenv.mkDerivation {
        name = "mediocre-blog-static-dep-shell";
        buildInputs = depInputs;
    };

    jekyllEnv = pkgs.bundlerEnv {
        name = "jekyllEnv";
        ruby = pkgs.ruby;
        gemdir = ./.;
    };

    build = pkgs.stdenv.mkDerivation {
        name = "mediocre-blog-static";
        src = ./src;
        buildPhase = "${jekyllEnv}/bin/jekyll build";
        installPhase = "mv _site $out";
    };

    serve = pkgs.writeScriptBin "static-serve" ''
        #!/bin/sh
        exec ${jekyllEnv}/bin/jekyll serve \
            -s ./src \
            -d ./_site \
            -w -I -D \
            -P 4002
    '';

    allInputs = depInputs ++ [ jekyllEnv serve ];

    shell = pkgs.stdenv.mkDerivation {
        name = "mediocre-blog-static-shell";
        buildInputs = allInputs;
    };
}
