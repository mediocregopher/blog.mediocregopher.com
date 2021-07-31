bundle:
	nix-build markov.nix -A nixBundle
	./result/bin/nix-bundle '((import ./markov.nix) {}).entrypoint' '/bin/markov'
