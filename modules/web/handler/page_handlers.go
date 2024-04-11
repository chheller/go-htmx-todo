package webhandlers

import (
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/domain"
	"github.com/chheller/go-htmx-todo/modules/user"
	"github.com/chheller/go-htmx-todo/modules/web"
)

type WebPageHandlers struct {
	UserService *user.UserService
}

// Entry point for /user routes. Configures handlers for the given router
func (wph WebPageHandlers) Init(router *http.ServeMux) domain.Handler {
	if wph.UserService == nil {
		panic("UserService must be initialized before calling InitializeHandlers")
	}
	router.HandleFunc("GET /signup", wph.handleGetUserPage)
	router.HandleFunc("GET /home", wph.handleGetHomePage)

	return wph
}

func (wph *WebPageHandlers) handleGetUserPage(res http.ResponseWriter, req *http.Request) {
	web.New().RenderTemplate(res, "user_signup.tmpl.html", &user.BasePageData{
		Title:   "Sign Up",
		DevMode: true,
	})
}

func (wph *WebPageHandlers) handleGetHomePage(res http.ResponseWriter, req *http.Request) {
	web.New().RenderTemplate(res, "base_page.tmpl.html", &user.BasePageData{
		DevMode: true,
	})
}
