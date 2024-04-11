package web

import (
	"embed"
	"html/template"
	"io"

	"github.com/sirupsen/logrus"
)

//go:embed "templates/*"
var templateFs embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	templates := template.Must(template.ParseFS(templateFs, "templates/*/*.html"))
	return &Template{templates: templates}
}

func (t *Template) RenderTemplate(w io.Writer, pathPrefix string, name string, data interface{}) error {
	tmpl := template.Must(t.templates.Clone())
	logrus.WithFields(logrus.Fields{"templates": tmpl.DefinedTemplates(), "page": "templates/" + name}).Debug("Executing template for page")

	tmpl = template.Must(tmpl.ParseFS(templateFs, "templates"+pathPrefix+"/"+name+".html"))

	return tmpl.ExecuteTemplate(w, name+".html", data)
}
