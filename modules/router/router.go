package router

import (
	"context"
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/types"
	"github.com/chheller/go-htmx-todo/modules/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateRouter creates a new http.ServeMux and initializes the http handlers.
// Wires up services with a database client, and attaches those services to handlers
// Responsible for middleware stack as well.
func CreateRouter(client *mongo.Client) *http.ServeMux {
	router := http.NewServeMux()
	handlers := []types.Handler{
		user.UserHandlers{
			UserService: (user.UserService{}).Init(client, context.Background()).(*user.UserService),
		},
	}

	for _, handler := range handlers {
		handler.Init(router)
	}

	return router
}
