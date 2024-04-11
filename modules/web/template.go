package web

import (
	"embed"
	"html/template"
	"io"
)

//go:embed "templates/*"
var templateFs embed.FS

type Template struct {
	templates *template.Template
}

var Templates = New()

func New() *Template {
	templates := template.Must(template.ParseFS(templateFs, "templates/*/*.html"))
	return &Template{templates: templates}
}

func (t *Template) RenderTemplate(w io.Writer, pathPrefix string, name string, data interface{}) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(templateFs, "templates"+pathPrefix+"/"+name+".html"))

	return tmpl.ExecuteTemplate(w, name+".html", data)
}
