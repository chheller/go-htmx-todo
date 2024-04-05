package user

import (
	"encoding/json"
	"net/http"
)

type UserHandlers struct {
	UserService *UserService
	Router      *http.ServeMux
}

func (uh *UserHandlers) InitializeHandlers() {
	if uh.UserService == nil || uh.Router == nil {
		panic("UserService and Router must be initialized before calling InitializeHandlers")
	}
	uh.Router.HandleFunc("GET /user", uh.handleGetUser)
	uh.Router.HandleFunc("POST /user", uh.handleCreateUser)
	uh.Router.HandleFunc("PUT /user", uh.handleUpdateUser)
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
