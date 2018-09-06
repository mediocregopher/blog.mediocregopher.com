# pacman -Sy ruby ruby-bundler
BUNDLE := bundle

serve:
	${BUNDLE} exec jekyll serve -w -I
	
update:
	${BUNDLE} update
