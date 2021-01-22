source $stdenv/setup
set -e

mkdir -p "$out"
$jekyll_env/bin/jekyll build -s "$src" -d "$out"
