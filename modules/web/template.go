package web

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"net/http"
)

//go:embed "templates/*"
var templateFs embed.FS

type Template struct {
	templates *template.Template
}

var Templates = New()

func New() *Template {
	templates := template.Must(template.ParseFS(templateFs, "templates/pages/base.layout.html"))
	return &Template{templates: templates}
}

func (t *Template) RenderTemplate(w io.Writer, pathPrefix string, name string, data interface{}) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.ParseFS(templateFs, "templates"+pathPrefix+"/"+name+".html"))

	return tmpl.ExecuteTemplate(w, name+".html", data)
}

func (t *Template) WriteTemplateResponse(w http.ResponseWriter, pathPrefix string, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	// Create a temporary buffer before writing out to the response writer,
	// so that if there is an error rendering the template,
	// we don't start writing to the response writer and then fail to complete the response.
	var temporaryWriter bytes.Buffer
	err := t.RenderTemplate(&temporaryWriter, pathPrefix, name, data)
	if err != nil {
		w.Header().Set("Location", "/500")
		w.WriteHeader(http.StatusMovedPermanently)
		// w.Write(nil)
		return
	}

	// Write the temporary buffer to the response writer if successful
	w.Write(temporaryWriter.Bytes())

}
