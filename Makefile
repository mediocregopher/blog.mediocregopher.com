serve:
	nix-shell -A serve

update:
	nix-shell -p bundler --run 'bundler update'

lock:
	nix-shell -p bundler -p bundix --run 'bundler lock; bundler package --no-install --path vendor; bundix; rm -rf .bundle vendor'

build:
	nix-build -A build
