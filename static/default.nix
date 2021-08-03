let
    utils = (import ../nix) {};
    pkgs = utils.pkgs;
    system = utils.system;

    jekyll_env = pkgs.bundlerEnv {
        name = "jekyll_env";
        ruby = pkgs.ruby;
        gemdir = ./.;
    };

    dep_inputs = [ pkgs.imagemagick pkgs.exiftool pkgs.bundler pkgs.bundix ];
    all_inputs = [ jekyll_env ] ++ dep_inputs;
in
    {
        build = derivation {
            inherit jekyll_env system;

            name = "mediocre-blog-static";
            builder = "${pkgs.bash}/bin/bash";
            args = [
                (pkgs.writeTextFile {
                    name = "mediocre-blog-static-buildsh";
                    text = ''
                        source ${pkgs.stdenv}/setup
                        set -e

                        mkdir -p "$out"
                        $jekyll_env/bin/jekyll build -s "${./src}" -d "$out"
                    '';
                    executable = true;
                })
            ];
        };

        dev = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-static-dev";
            buildInputs = all_inputs;
            shellHook = ''
                exec ${jekyll_env}/bin/jekyll serve -s ./src -d ./_site -w -I -D -H 0.0.0.0 -P 4001
            '';
        };

        depShell = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-static-dep-shell";
            buildInputs = dep_inputs;
        };

        shell = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-static-shell";
            buildInputs = all_inputs;
        };
    }
