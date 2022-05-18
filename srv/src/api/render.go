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

func mustReadTplFile(fileName string) string {
	path := filepath.Join("tpl", fileName)

	b, err := fs.ReadFile(tplFS, path)
	if err != nil {
		panic(fmt.Errorf("reading file %q from tplFS: %w", path, err))
	}

	return string(b)
}

func (a *api) mustParseTpl(name string) *template.Template {

	blogURL := func(path string) string {

		trailingSlash := strings.HasSuffix(path, "/")
		path = filepath.Join(a.params.PathPrefix, "/v2", path)

		if trailingSlash {
			path += "/"
		}

		return path
	}

	tpl := template.New("").Funcs(template.FuncMap{
		"BlogURL": blogURL,
		"AssetURL": func(path string) string {
			path = filepath.Join("assets", path)
			return blogURL(path)
		},
	})

	tpl = template.Must(tpl.Parse(mustReadTplFile(name)))

	return tpl
}

func (a *api) mustParseBasedTpl(name string) *template.Template {
	tpl := a.mustParseTpl(name)
	tpl = template.Must(tpl.New("base.html").Parse(mustReadTplFile("base.html")))
	return tpl
}

type tplData struct {
	Payload   interface{}
	CSRFToken string
}

func (t tplData) CSRFFormInput() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<input type="hidden" name="%s" value="%s" />`,
		csrfTokenFormName, t.CSRFToken,
	))
}

// executeTemplate expects to be the final action in an http.Handler
func executeTemplate(
	rw http.ResponseWriter, r *http.Request,
	tpl *template.Template, payload interface{},
) {

	csrfToken, _ := apiutil.GetCookie(r, csrfTokenCookieName, "")

	tplData := tplData{
		Payload:   payload,
		CSRFToken: csrfToken,
	}

	if err := tpl.Execute(rw, tplData); err != nil {
		apiutil.InternalServerError(
			rw, r, fmt.Errorf("rendering template: %w", err),
		)
		return
	}
}

func (a *api) executeRedirectTpl(
	rw http.ResponseWriter, r *http.Request, path string,
) {
	executeTemplate(rw, r, a.redirectTpl, struct {
		Path string
	}{
		Path: path,
	})
}

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

		posts, hasMore, err := a.params.PostStore.WithOrderDesc().Get(page, pageCount)
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

func (a *api) renderDumbHandler(tplName string) http.Handler {

	tpl := a.mustParseBasedTpl(tplName)

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := tpl.Execute(rw, nil); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("rendering %q: %w", tplName, err),
			)
			return
		}
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
