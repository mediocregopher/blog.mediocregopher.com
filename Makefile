# pacman -Sy ruby ruby-bundler
BUNDLE := bundle

serve:
	${BUNDLE} exec jekyll serve -w -I -D -H 0.0.0.0
	
update:
	${BUNDLE} update
