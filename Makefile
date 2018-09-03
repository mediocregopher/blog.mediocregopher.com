# pacman -Sy ruby ruby-bundler
BUNDLE := bundle

serve:
	${BUNDLE} exec jekyll serve -w
	
update:
	${BUNDLE} update
