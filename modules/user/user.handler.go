package user

import (
	"encoding/json"
	"net/http"
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
