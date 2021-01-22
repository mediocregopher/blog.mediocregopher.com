{
    pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/cd63096d6d887d689543a0b97743d28995bc9bc3.tar.gz") {},
    system ? builtins.currentSystem,
}:

    let
        jekyll_env = pkgs.bundlerEnv {
            name = "jekyll_env";
            ruby = pkgs.ruby;
            gemdir = ./.;
        };
    in
        {
            build = derivation {
                system = system;
                name = "mediocre-blog";
                builder = "${pkgs.bash}/bin/bash";
                args = [ ./build.sh ];

                src = ./src;
                stdenv = pkgs.stdenv;
                inherit jekyll_env;
            };

            serve = pkgs.stdenv.mkDerivation rec {
                name = "jekyll_env";

                # glibcLocales is required so to fill in LC_ALL and other locale
                # related environment vars. Without those jekyll's scss compiler
                # fails.
                #
                # TODO probably get rid of the scss compiler.
                buildInputs = [ jekyll_env pkgs.glibcLocales ];

                shellHook = ''
                    exec ${jekyll_env}/bin/jekyll serve -s ./src -d ./_site -w -I -D -H 0.0.0.0
                '';
            };
        }


