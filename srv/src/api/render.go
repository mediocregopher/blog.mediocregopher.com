package api

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

//go:embed tpl
var tplFS embed.FS

func mustParseTpl(name string) *template.Template {

	mustRead := func(fileName string) string {
		path := filepath.Join("tpl", fileName)

		b, err := fs.ReadFile(tplFS, path)
		if err != nil {
			panic(fmt.Errorf("reading file %q from tplFS: %w", path, err))
		}

		return string(b)
	}

	tpl := template.Must(template.New("").Parse(mustRead(name)))
	tpl = template.Must(tpl.New("base.html").Parse(mustRead("base.html")))

	return tpl
}

func (a *api) renderIndexHandler() http.Handler {

	tpl := mustParseTpl("index.html")
	const pageCount = 20

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

		posts, _, err := a.params.PostStore.Get(page, pageCount)
		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching page %d of posts: %w", page, err),
			)
			return
		}

		tplData := struct {
			Posts []post.StoredPost
		}{
			Posts: posts,
		}

		if err := tpl.Execute(rw, tplData); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering index: %w", err),
			)
			return
		}
	})
}

func (a *api) renderPostHandler() http.Handler {

	tpl := mustParseTpl("post.html")

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

		if err := tpl.Execute(rw, tplData); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering post with id %q: %w", id, err),
			)
			return
		}
	})
}
