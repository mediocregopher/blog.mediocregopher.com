package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

func (a *api) renderIndexHandler() http.Handler {

	tpl := a.mustParseBasedTpl("index.html")
	const pageCount = 10

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if path := r.URL.Path; !strings.HasSuffix(path, "/") && filepath.Base(path) != "index.html" {
			http.Error(rw, "Page not found", 404)
			return
		}

		page, err := apiutil.StrToInt(r.FormValue("p"), 0)
		if err != nil {
			apiutil.BadRequest(
				rw, r, fmt.Errorf("invalid page number: %w", err),
			)
			return
		}

		posts, hasMore, err := a.params.PostStore.Get(page, pageCount)
		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching page %d of posts: %w", page, err),
			)
			return
		}

		tplPayload := struct {
			Posts              []post.StoredPost
			PrevPage, NextPage int
		}{
			Posts:    posts,
			PrevPage: -1,
			NextPage: -1,
		}

		if page > 0 {
			tplPayload.PrevPage = page - 1
		}

		if hasMore {
			tplPayload.NextPage = page + 1
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}
