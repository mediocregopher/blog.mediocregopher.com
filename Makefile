result:
	nix-build -A build

install: result
	nix-env -i "$$(readlink result)"

clean:
	rm -f result
	rm -rf _site

serve:
	nix-shell -A serve

update:
	nix-shell -p bundler --run 'bundler update; bundler lock; bundix; rm -rf .bundle vendor'
