# pacman -Sy ruby ruby-bundler
BUNDLE := bundle

serve:
	${BUNDLE} exec jekyll serve -w -I -D
	
update:
	${BUNDLE} update
