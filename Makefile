result:
	nix-build -A build

install: result
	nix-env -i "$$(readlink result)"

clean:
	rm -f result
	rm -rf _site

serve:
	nix-shell -A serve

lock:
	nix-shell -p bundler -p bundix --run 'bundler lock; bundix; rm -rf .bundle vendor'

update:
	nix-shell -p bundler -p bundix --run 'bundler update; bundler lock; bundix; rm -rf .bundle vendor'
