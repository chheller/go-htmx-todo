package webhandlers

import (
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/domain"
	"github.com/chheller/go-htmx-todo/modules/user"
	"github.com/chheller/go-htmx-todo/modules/web"
	viewmodel "github.com/chheller/go-htmx-todo/modules/web/view_model"
	"github.com/sirupsen/logrus"
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
	err := web.Templates.RenderTemplate(res, "/pages", "user_signup", viewmodel.DefaultSignupPageData)
	if err != nil {
		logrus.WithError(err).Error("Failed to render template")
	}
}

func (wph *WebPageHandlers) handleGetHomePage(res http.ResponseWriter, req *http.Request) {
	err := web.Templates.RenderTemplate(res, "/pages", "base_page", viewmodel.DefaultBasePageData)
	if err != nil {
		logrus.WithError(err).Error("Failed to render template")

	}

}
