let
    utils = (import ../nix) {};
    pkgs = utils.pkgs;
    system = utils.system;
in
    {

        build = pkgs.buildGoModule {
            pname = "mediocre-blog-srv";
            version = "dev";
            src = ./.;
            vendorSha256 = "08wv94yv2wmlxzmanw551gixc8v8nl6zq2m721ig9nl3r540x46f";
        };

        shell = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-srv-shell";
            buildInputs = [ pkgs.go ];
        };

    }
