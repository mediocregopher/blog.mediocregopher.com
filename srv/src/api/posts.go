package api

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutils"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

func (a *api) postHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := strings.TrimSuffix(filepath.Base(r.URL.Path), ".html")

		storedPost, err := a.params.PostStore.GetByID(id)

		if errors.Is(err, post.ErrPostNotFound) {
			http.Error(rw, "Post not found", 404)
			return
		} else if err != nil {
			apiutils.InternalServerError(
				rw, r, fmt.Errorf("fetching post with id %q: %w", id, err),
			)
			return
		}

		renderablePost, err := post.NewRenderablePost(a.params.PostStore, storedPost)
		if err != nil {
			apiutils.InternalServerError(
				rw, r, fmt.Errorf("constructing renderable post with id %q: %w", id, err),
			)
			return
		}

		if err := a.params.PostHTTPRenderer.Render(rw, renderablePost); err != nil {
			apiutils.InternalServerError(
				rw, r, fmt.Errorf("rendering post with id %q: %w", id, err),
			)
			return
		}
	})
}
