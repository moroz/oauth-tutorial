package templates

import (
	"embed"
	"html/template"
)

//go:embed views/*.html.tmpl
var fs embed.FS

var PagesIndex = template.Must(template.ParseFS(fs, "views/layout.html.tmpl", "views/index.html.tmpl"))
