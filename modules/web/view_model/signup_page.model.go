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

type SigninPageData struct {
	*BasePageData
	VerificationToken string
}

var DefaultSigninPageData = func (verificationToken string) *SigninPageData {
	return &SigninPageData{
		BasePageData: &BasePageData{
			Title:                     "Signing In",
			InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
		},
		VerificationToken: verificationToken,
	}
}

type SignupCompletePageData struct {
	*BasePageData
}

var DefaultSignupCompletePageData = &SignupCompletePageData{
	BasePageData: &BasePageData{
		Title:                     "Signup",
		InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,

	},
}