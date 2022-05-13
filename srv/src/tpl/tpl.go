// Package tpl contains template files which are used to render the blog.
package tpl

import (
	"embed"
	html_tpl "html/template"
)

//go:embed *
var fs embed.FS

var HTML = html_tpl.Must(html_tpl.ParseFS(fs, "html/*"))
