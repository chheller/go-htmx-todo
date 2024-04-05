package router

import (
	"context"
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/user"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateRouter(client *mongo.Client) *http.ServeMux {
	router := http.NewServeMux()
	userHandler := user.UserHandlers{
		UserService: user.UserService{
			Client: client,
			Ctx:    context.Background(),
		},
		Router: router,
	}
	userHandler.InitializeHandlers()

	return router
}
