
test_dir="$(mktemp -d)"

mkdir -p "$test_dir"/run
mkdir -p "$test_dir"/data

test_cfg="(import ./config.nix) // {
    runDir=\"${test_dir}/run\";
    dataDir=\"${test_dir}/data\";
}"

entrypoint=$(nix-build --no-out-link -A entrypoint \
    --arg baseConfig "$test_cfg" \
    --arg skipServices '["srv"]')

$entrypoint &
trap "kill $!; wait; rm -rf $test_dir" EXIT

# NOTE this is a bit of a hack... the location of the redis socket's source of
# truth is in default.nix, but it's not clear how to get that from there to
# here, so we reproduce the calculation here.
while [ ! -e $test_dir/run/redis ]; do
    echo "waiting for redis unix socket"
    sleep 1
done

nix-shell -A srv.shell \
    --arg baseConfig "$test_cfg" \
    "$@"
