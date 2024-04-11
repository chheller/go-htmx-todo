package web

import (
	"embed"
	"html/template"
	"io"

	"github.com/sirupsen/logrus"
)

//go:embed templates/base_page.tmpl.html templates/**/*.tmpl*
var templateFs embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	templates := template.Must(template.ParseFS(templateFs, "templates/pages/*.tmpl*", "templates/components/*.tmpl*", "templates/email/*.tmpl*", "templates/base_page.tmpl.html"))
	return &Template{templates: templates}
}

func (t *Template) RenderTemplate(w io.Writer, name string, data interface{}) error {
	tmpl := template.Must(t.templates.Clone())
	logrus.WithField("templates", tmpl.DefinedTemplates()).Debug("Loaded templates")

	tmpl = template.Must(tmpl.ParseFS(templateFs, "templates/**/"+name))

	return tmpl.ExecuteTemplate(w, name, data)
}
