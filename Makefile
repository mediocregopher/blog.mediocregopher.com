BUNDLE := ~/.gem/ruby/2.5.0/bin/bundle

serve:
	${BUNDLE} exec jekyll serve
	
update:
	${BUNDLE} update
