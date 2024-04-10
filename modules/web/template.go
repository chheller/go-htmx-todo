package web

import (
	"embed"
	"text/template"

	"github.com/sirupsen/logrus"
)

//go:embed templates/**/*
var templateFs embed.FS
var templates *template.Template

func parseTemplates() (*template.Template, error) {
	return template.ParseFS(templateFs, "templates/*/*.go.tmpl")
}

func GetTemplates() *template.Template {
	if templates == nil {
		var err error
		templates, err = parseTemplates()
		if err != nil {
			logrus.WithError(err).Panic("Failed to parse templates")
		}
	}
	return templates
}
