package user

import (
	"encoding/json"
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/web"
	viewmodel "github.com/chheller/go-htmx-todo/modules/web/view_model"
	log "github.com/sirupsen/logrus"
)

type UserHandlers struct {
	UserService *UserService
}

// Entry point for /user routes. Configures handlers for the given router
func (uh UserHandlers) Init(router *http.ServeMux) {
	if uh.UserService == nil {
		panic("UserService must be initialized before calling InitializeHandlers")
	}
	router.HandleFunc("GET /signup", uh.handleGetUserPage)
	router.HandleFunc("POST /signup", uh.handleCreateUser)
	router.HandleFunc("GET /home", uh.handleGetHomePage)
	router.HandleFunc("GET /signup/verify", uh.handleVerifyUserOtp)
}

func (u *UserHandlers) handleCreateUser(res http.ResponseWriter, req *http.Request) {
	user := &User{}
	json.NewDecoder(req.Body).Decode(user)
	err := u.UserService.CreateUser(*user, req.Context())
	if err != nil {
		log.WithError(err).Error("Failed to create user")
		// res.WriteHeader(http.StatusInternalServerError)
		web.Templates.WriteTemplateResponse(res, "/components/user", "create_user_error", nil)
		return
	}
	web.Templates.WriteTemplateResponse(res, "/components/user", "create_user_success", nil )
}

func (u *UserHandlers) handleVerifyUserOtp(res http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		log.Info("token param not found in query params")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Invalid token"))
		return
	}

	if ok := u.UserService.VerifyUserOtp(token, req.Context()); !ok {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Bad Authorization"))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))

}
func (u *UserHandlers) handleGetUserPage(res http.ResponseWriter, req *http.Request) {
	web.Templates.WriteTemplateResponse(res, "/pages/user", "user_signup", viewmodel.DefaultSignupPageData)
}

func (u *UserHandlers) handleGetHomePage(res http.ResponseWriter, req *http.Request) {
	web.Templates.WriteTemplateResponse(res, "/pages", "base_page", viewmodel.DefaultBasePageData)
}
