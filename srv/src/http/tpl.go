package http

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
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

func (a *api) parseTpl(tplBody string) (*template.Template, error) {

	blogURL := func(path string) string {

		// filepath.Join strips trailing slash, but we want to keep it
		trailingSlash := strings.HasSuffix(path, "/")

		path = filepath.Join("/", path)

		if trailingSlash && path != "/" {
			path += "/"
		}

		return path
	}

	tpl := template.New("root")

	tpl = tpl.Funcs(template.FuncMap{
		"BlogURL": blogURL,
		"StaticURL": func(path string) string {
			path = filepath.Join("static", path)
			return blogURL(path)
		},
		"AssetURL": func(id string) string {
			path := filepath.Join("assets", id)
			return blogURL(path)
		},
		"PostURL": func(id string) string {
			path := filepath.Join("posts", id)
			return blogURL(path)
		},
		"DateTimeFormat": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	})

	tpl = template.Must(tpl.New("image.html").Parse(mustReadTplFile("image.html")))

	tpl = tpl.Funcs(template.FuncMap{
		"Image": func(id string) (template.HTML, error) {

			tplPayload := struct {
				ID        string
				Resizable bool
			}{
				ID:        id,
				Resizable: isImgResizable(id),
			}

			buf := new(bytes.Buffer)
			if err := tpl.ExecuteTemplate(buf, "image.html", tplPayload); err != nil {
				return "", err
			}

			return template.HTML(buf.Bytes()), nil
		},
	})

	var err error

	if tpl, err = tpl.New("").Parse(tplBody); err != nil {
		return nil, err
	}

	return tpl, nil
}

func (a *api) mustParseTpl(name string) *template.Template {
	return template.Must(a.parseTpl(mustReadTplFile(name)))
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
