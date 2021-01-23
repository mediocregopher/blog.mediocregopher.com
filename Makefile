result:
	nix-build -A build

install: result
	nix-env -i "$$(readlink result)"

clean:
	rm result

serve:
	nix-shell -A serve

update:
	nix-shell -p bundler --run 'bundler update'

lock:
	nix-shell -p bundler -p bundix --run 'bundler lock; bundler package --no-install --path vendor; bundix; rm -rf .bundle vendor'
