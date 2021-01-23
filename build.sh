source $stdenv/setup
set -e

d="$out/var/www/blog.mediocregopher.com"
mkdir -p "$d"
$jekyll_env/bin/jekyll build -s "$src" -d "$d"
