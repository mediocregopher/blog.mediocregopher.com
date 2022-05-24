package http

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
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

func (a *api) blogURL(path string, abs bool) string {
	// filepath.Join strips trailing slash, but we want to keep it
	trailingSlash := strings.HasSuffix(path, "/")

	res := filepath.Join("/", path)

	if trailingSlash && res != "/" {
		res += "/"
	}

	if abs {
		res = a.params.PublicURL.String() + res
	}

	return res
}

func (a *api) postURL(id string, abs bool) string {
	path := filepath.Join("posts", id)
	return a.blogURL(path, abs)
}

func (a *api) postsURL(abs bool) string {
	return a.blogURL("posts", abs)
}

func (a *api) assetsURL(abs bool) string {
	return a.blogURL("assets", abs)
}

func (a *api) tplFuncs() template.FuncMap {
	return template.FuncMap{
		"BlogURL": func(path string) string {
			return a.blogURL(path, false)
		},
		"StaticURL": func(path string) string {
			path = filepath.Join("static", path)
			return a.blogURL(path, false)
		},
		"AssetURL": func(id string) string {
			path := filepath.Join("assets", id)
			return a.blogURL(path, false)
		},
		"PostURL": func(id string) string {
			return a.postURL(id, false)
		},
		"DateTimeFormat": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	}
}

func (a *api) parseTpl(name, tplBody string) (*template.Template, error) {

	tpl := template.New(name)
	tpl = tpl.Funcs(a.tplFuncs())

	var err error

	if tpl, err = tpl.Parse(tplBody); err != nil {
		return nil, err
	}

	return tpl, nil
}

func (a *api) mustParseTpl(name string) *template.Template {
	return template.Must(a.parseTpl(name, mustReadTplFile(name)))
}

func (a *api) mustParseBasedTpl(name string) *template.Template {
	tpl := a.mustParseTpl(name)
	tpl = template.Must(tpl.New("load-csrf.html").Parse(mustReadTplFile("load-csrf.html")))
	tpl = template.Must(tpl.New("base.html").Parse(mustReadTplFile("base.html")))
	return tpl
}

type tplData struct {
	Payload   interface{}
	CSRFToken string
}

func (t tplData) CSRFFormInput() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<input type="hidden" name="%s" class="csrfHiddenInput" />`,
		csrfTokenFormName,
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
	rw http.ResponseWriter, r *http.Request, url string,
) {
	log.Printf("here url:%q", url)
	executeTemplate(rw, r, a.redirectTpl, struct {
		URL string
	}{
		URL: url,
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
