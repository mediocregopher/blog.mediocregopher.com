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

        shell_inputs = [ jekyll_env pkgs.imagemagick pkgs.exiftool ];
    in
        {
            build = derivation {
                inherit jekyll_env system;

                name = "mediocre-blog";
                builder = "${pkgs.bash}/bin/bash";
                args = [ ./build.sh ];

                src = ./src;
                stdenv = pkgs.stdenv;
            };

            serve = pkgs.stdenv.mkDerivation {
                name = "mediocre-blog-shell-serve";
                buildInputs = shell_inputs;
                shellHook = ''
                    exec ${jekyll_env}/bin/jekyll serve -s ./src -d ./_site -w -I -D -H 0.0.0.0
                '';
            };

            shell = pkgs.stdenv.mkDerivation {
                name = "mediocre-blog-shell";
                buildInputs = shell_inputs;
            };
        }
