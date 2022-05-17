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

func (a *api) mustParseTpl(name string) *template.Template {

	mustRead := func(fileName string) string {
		path := filepath.Join("tpl", fileName)

		b, err := fs.ReadFile(tplFS, path)
		if err != nil {
			panic(fmt.Errorf("reading file %q from tplFS: %w", path, err))
		}

		return string(b)
	}

	blogURL := func(path string) string {
		return filepath.Join(a.params.PathPrefix, "/v2", path)
	}

	tpl := template.New("").Funcs(template.FuncMap{
		"BlogURL": blogURL,
		"AssetURL": func(path string) string {
			path = filepath.Join("assets", path)
			return blogURL(path)
		},
	})

	tpl = template.Must(tpl.Parse(mustRead(name)))
	tpl = template.Must(tpl.New("base.html").Parse(mustRead("base.html")))

	return tpl
}

func (a *api) renderIndexHandler() http.Handler {

	tpl := a.mustParseTpl("index.html")
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

		posts, hasMore, err := a.params.PostStore.WithOrderDesc().Get(page, pageCount)
		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching page %d of posts: %w", page, err),
			)
			return
		}

		tplData := struct {
			Posts              []post.StoredPost
			PrevPage, NextPage int
		}{
			Posts:    posts,
			PrevPage: -1,
			NextPage: -1,
		}

		if page > 0 {
			tplData.PrevPage = page - 1
		}

		if hasMore {
			tplData.NextPage = page + 1
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

	tpl := a.mustParseTpl("post.html")

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

func (a *api) renderDumbHandler(tplName string) http.Handler {

	tpl := a.mustParseTpl(tplName)

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := tpl.Execute(rw, nil); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering %q: %w", tplName, err),
			)
			return
		}
	})
}

func (a *api) renderAdminAssets() http.Handler {

	tpl := a.mustParseTpl("admin/assets.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ids, err := a.params.PostAssetStore.List()

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("getting list of asset ids: %w", err),
			)
			return
		}

		tplData := struct {
			IDs []string
		}{
			IDs: ids,
		}

		if err := tpl.Execute(rw, tplData); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering: %w", err),
			)
			return
		}
	})
}
