
CONFIG = ./config.nix

entrypoint:
	nix-build -A entrypoint \
		--arg baseConfig '(import ${CONFIG})'

install:
	$$(nix-build -A install --arg baseConfig '(import ${CONFIG})')

test:
	$$(nix-build --no-out-link -A pkgs.bash)/bin/bash srv-dev-env.sh \
    	--run "cd srv/src && go test ./... -count=1 -tags integration"
	@echo "\nTESTS PASSED!\n"

srv.dev-shell:
	$$(nix-build --no-out-link -A pkgs.bash)/bin/bash srv-dev-env.sh \
		--command " \
			cd srv/src; \
			go run cmd/import-posts/main.go ../../static/src/_posts/*; \
			return; \
		"

srv.shell:
	nix-shell -A srv.shellWithBuild --arg baseConfig '(import ${CONFIG})' \
		--command 'cd srv/src; return'

# TODO static is on the way out, these aren't well supported
static.serve:
	nix-shell -A static.shell --run 'cd static; static-serve'

static.depShell:
	nix-shell -A static.depShell --command 'cd static; return'

static.lock:
	nix-shell -A static.depShell  --run 'bundler lock; bundix; rm -rf .bundle vendor'
