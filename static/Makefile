build:
	nix-build -A build

clean:
	rm -f result
	rm -rf _site

dev:
	nix-shell -A dev

shell:
	nix-shell -A shell

lock:
	nix-shell -A depShell --run 'bundler lock; bundix; rm -rf .bundle vendor'

update:
	nix-shell -A depShell --run 'bundler update; bundler lock; bundix; rm -rf .bundle vendor'
