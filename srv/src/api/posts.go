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

func (a *api) postHandler() http.Handler {
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

		tplData := struct {
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
					tplData.SeriesPrevious = &seriesPost
					continue
				}

				tplData.SeriesNext = &seriesPost
				break
			}
		}

		if err := tpls.ExecuteTemplate(rw, "post.html", tplData); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering post with id %q: %w", id, err),
			)
			return
		}
	})
}
