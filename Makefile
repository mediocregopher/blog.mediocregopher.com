
CONFIG = ./config.nix

entrypoint:
	nix-build -A entrypoint \
		--arg baseConfig '(import ${CONFIG})'

install:
	nix-build -A install --arg baseConfig '(import ${CONFIG})'
	./result

test:
	$$(nix-build --no-out-link -A pkgs.bash)/bin/bash test.sh
	@if [ $$? == 0 ]; then echo "TESTS PASSED!"; else echo "TESTS FAILED!"; fi

srv.shell:
	nix-shell -A srv.shell --arg baseConfig '(import ${CONFIG})' \
		--command 'cd srv; return'

# TODO static is on the way out, these aren't well supported
static.serve:
	nix-shell -A static.shell --run 'cd static; static-serve'

static.depShell:
	nix-shell -A static.depShell --command 'cd static; return'

static.lock:
	nix-shell -A static.depShell  --run 'bundler lock; bundix; rm -rf .bundle vendor'
