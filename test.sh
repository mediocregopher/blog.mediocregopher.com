
test_dir="$(mktemp -d)"

mkdir -p "$test_dir"/run
mkdir -p "$test_dir"/data

test_cfg="(import ./config.nix) // {
    runDir=\"${test_dir}/run\";
    dataDir=\"${test_dir}/data\";
}"

$(nix-build --no-out-link -A entrypoint \
    --arg baseConfig "$test_cfg" \
    --arg baseSkipServices '["srv" "static"]') &

trap "kill $!; wait; rm -rf $test_dir" EXIT

# TODO there's a race condition here, we should wait until redis is definitely
# listening before commencing the tests.

nix-shell -A srv.test \
    --arg baseConfig "$test_cfg"  \
    --run "cd srv/src && go test ./... -count=1 -tags integration"
