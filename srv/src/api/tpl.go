package api

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutil"
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

func (a *api) renderDumbTplHandler(tplName string) http.Handler {

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
