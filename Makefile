serve:
	docker run -it --rm \
		-v $$(pwd):/srv/jekyll \
		-p 4000:4000 \
		jekyll/jekyll \
			jekyll serve -w -I -D -H 0.0.0.0

update:
	docker run -it --rm \
		-v $$(pwd):/srv/jekyll \
		jekyll/jekyll \
			bundle update
