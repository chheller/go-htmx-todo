package viewmodel

import "github.com/chheller/go-htmx-todo/modules/config"

// TODO: Move this and initialize it for every page to pull in
type BasePageData struct {
	Title                     string
	InjectBrowserReloadScript bool
}

var DefaultBasePageData = &BasePageData{
	Title:                     "Go HTMX Todo",
	InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
}
