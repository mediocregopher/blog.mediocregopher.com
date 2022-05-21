package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

func (a *api) renderIndexHandler() http.Handler {

	legacyPostPathRegexp := regexp.MustCompile(
		`^/[0-9]{4}/[0-9]{2}/[0-9]{2}/([^/.]+)\.html$`,
	)

	tpl := a.mustParseBasedTpl("index.html")
	const pageCount = 10

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		if matches := legacyPostPathRegexp.FindStringSubmatch(path); len(matches) == 2 {
			id := matches[1]
			http.Redirect(rw, r, filepath.Join("/posts", id), http.StatusMovedPermanently)
			return
		}

		if !strings.HasSuffix(path, "/") && filepath.Base(path) != "index.html" {
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

		tag := r.FormValue("tag")

		var (
			posts   []post.StoredPost
			hasMore bool
		)

		if tag == "" {
			posts, hasMore, err = a.params.PostStore.Get(page, pageCount)
		} else {
			posts, err = a.params.PostStore.GetByTag(tag)
		}

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching page %d of posts: %w", page, err),
			)
			return
		}

		tags, err := a.params.PostStore.GetTags()
		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fething tags: %w", err),
			)
			return
		}

		tplPayload := struct {
			Posts              []post.StoredPost
			PrevPage, NextPage int
			Tags               []string
		}{
			Posts:    posts,
			PrevPage: -1,
			NextPage: -1,
			Tags:     tags,
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
