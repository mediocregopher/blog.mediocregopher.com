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
            vendorSha256 = "1l8783zxiv8p74xr5y61s6r2j3mrrgx714i546l6cy0qfjhk7s7m";
        };

        shell = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-srv-shell";
            buildInputs = [ pkgs.go ];
        };

    }
