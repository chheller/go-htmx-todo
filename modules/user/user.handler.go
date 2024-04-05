package user

import (
	"encoding/json"
	"net/http"
)

type UserHandlers struct {
	UserService UserService
	Router      *http.ServeMux
}

func (uh *UserHandlers) InitializeHandlers() {

	handleGetUser := func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("GET /user"))
	}

	handleUpdateUser := func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("PUT /user"))
	}

	uh.Router.HandleFunc("GET /user", handleGetUser)
	uh.Router.HandleFunc("POST /user", handleCreateUser(&uh.UserService))
	uh.Router.HandleFunc("PUT /user", handleUpdateUser)

}

func handleCreateUser(svc *UserService) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		user := &User{}
		json.NewDecoder(req.Body).Decode(user)
		svc.CreateUser(*user)
		res.Write([]byte("POST /user"))
	}
}
