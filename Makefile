SKIP_SERVICES = []

all:
	nix-build -A entrypoint \
		--arg baseConfig '(import ./config.nix)' \
		--arg baseSkipServices '${SKIP_SERVICES}'

run: all
	./result

all.prod:
	nix-build -A entrypoint \
		--arg baseConfig '(import ./prod.config.nix)' \
		--arg baseSkipServices '${SKIP_SERVICES}'

run.prod: all.prod
	./result

install.prod:
	nix-build -A install --arg baseConfig '(import ./prod.config.nix)'
	./result

srv.shell:
	nix-shell -A srv.shell --command 'cd srv; return'

srv.shell.prod:
	nix-shell -A srv.shell --arg baseConfig '(import ./prod.config.nix)' --command 'cd srv; return'

static.shell:
	nix-shell -A static.shell --command 'cd static; return'

static.serve:
	nix-shell -A static.shell --run 'cd static; static-serve'

static.depShell:
	nix-shell -A static.depShell --command 'cd static; return'

static.lock:
	nix-shell -A static.depShell  --run 'bundler lock; bundix; rm -rf .bundle vendor'
