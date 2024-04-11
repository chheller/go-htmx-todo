package viewmodel

import "github.com/chheller/go-htmx-todo/modules/config"

type SignupPageData struct {
	BasePageData *BasePageData
}

var DefaultSignupPageData = &SignupPageData{
	BasePageData: &BasePageData{
		Title:                     "Signup",
		InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
	},
}
