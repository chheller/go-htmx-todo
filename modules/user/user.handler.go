package user

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type UserHandlers struct {
	UserService *UserService
}

// Entry point for /user routes. Configures handlers for the given router
func (uh *UserHandlers) InitializeHandlers(router *http.ServeMux) {
	if uh.UserService == nil {
		panic("UserService must be initialized before calling InitializeHandlers")
	}
	router.HandleFunc("GET /user", uh.handleGetUser)
	router.HandleFunc("POST /user", uh.handleCreateUser)
	router.HandleFunc("PUT /user", uh.handleUpdateUser)
	router.HandleFunc("GET /user/verify-otp", uh.handleVerifyUserOtp)
}

func (u *UserHandlers) handleCreateUser(res http.ResponseWriter, req *http.Request) {
	user := &User{}
	json.NewDecoder(req.Body).Decode(user)
	u.UserService.CreateUser(*user)
	res.Write([]byte("POST /user"))
}
func (u *UserHandlers) handleUpdateUser(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("PUT /user"))
}

func (u *UserHandlers) handleGetUser(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("GET /user"))
}

func (u *UserHandlers) handleVerifyUserOtp(res http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		log.Print("token param not found in query params")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Invalid token"))
		return
	}

	if ok := u.UserService.VerifyUserOtp(token); !ok {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte("Bad Authorization"))
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))

}
