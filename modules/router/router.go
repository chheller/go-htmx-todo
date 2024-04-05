package router

import (
	"context"
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateRouter creates a new http.ServeMux and initializes the http handlers.
// Wires up services with a database client, and attaches those services to handlers
// Responsible for middleware stack as well.
func CreateRouter(client *mongo.Client) *http.ServeMux {
	router := http.NewServeMux()
	userHandler := user.UserHandlers{
		UserService: &user.UserService{
			Client: client,
			Ctx:    context.Background(),
		},
	}
	userHandler.InitializeHandlers(router)

	return router
}
