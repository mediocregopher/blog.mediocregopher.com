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
            vendorSha256 = "0xr5gks5mrh34s5npncw71wncrzqrhnm3vjfwdakd7fzd6iw049z";
        };

        shell = pkgs.stdenv.mkDerivation {
            name = "mediocre-blog-srv-shell";
            buildInputs = [ pkgs.go ];
        };

    }
