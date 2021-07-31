{
    pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/d50923ab2d308a1ddb21594ba6ae064cab65d8ae.tar.gz") {},
    system ? builtins.currentSystem,
}:
    {
        pkgs = pkgs;
        system = system;
    }

