package api

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

func (a *api) renderPostHandler() http.Handler {

	tpl := a.mustParseBasedTpl("post.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := strings.TrimSuffix(filepath.Base(r.URL.Path), ".html")

		storedPost, err := a.params.PostStore.GetByID(id)

		if errors.Is(err, post.ErrPostNotFound) {
			http.Error(rw, "Post not found", 404)
			return
		} else if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching post with id %q: %w", id, err),
			)
			return
		}

		parserExt := parser.CommonExtensions | parser.AutoHeadingIDs
		parser := parser.NewWithExtensions(parserExt)

		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

		renderedBody := markdown.ToHTML([]byte(storedPost.Body), parser, htmlRenderer)

		tplPayload := struct {
			post.StoredPost
			SeriesPrevious, SeriesNext *post.StoredPost
			Body                       template.HTML
		}{
			StoredPost: storedPost,
			Body:       template.HTML(renderedBody),
		}

		if series := storedPost.Series; series != "" {

			seriesPosts, err := a.params.PostStore.GetBySeries(series)
			if err != nil {
				apiutil.InternalServerError(
					rw, r,
					fmt.Errorf("fetching posts for series %q: %w", series, err),
				)
				return
			}

			var foundThis bool

			for i := range seriesPosts {

				seriesPost := seriesPosts[i]

				if seriesPost.ID == storedPost.ID {
					foundThis = true
					continue
				}

				if !foundThis {
					tplPayload.SeriesPrevious = &seriesPost
					continue
				}

				tplPayload.SeriesNext = &seriesPost
				break
			}
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}

func (a *api) renderPostAssetsIndexHandler() http.Handler {

	tpl := a.mustParseBasedTpl("assets.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ids, err := a.params.PostAssetStore.List()

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("getting list of asset ids: %w", err),
			)
			return
		}

		tplPayload := struct {
			IDs []string
		}{
			IDs: ids,
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}
