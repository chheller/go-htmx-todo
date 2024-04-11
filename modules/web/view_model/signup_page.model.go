package viewmodel

import "github.com/chheller/go-htmx-todo/modules/config"

type SignupPageData struct {
	*BasePageData
}

var DefaultSignupPageData = &SignupPageData{
	BasePageData: &BasePageData{
		Title:                     "Sign Up",
		InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
	},
}
