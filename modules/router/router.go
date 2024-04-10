package router

import (
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/domain"
	"github.com/chheller/go-htmx-todo/modules/user"
)

// CreateRouter creates a new http.ServeMux and initializes the http handlers.
// Wires up services with a database client, and attaches those services to handlers
// Responsible for middleware stack as well.
func CreateRouter(services *ApplicationServices) *http.ServeMux {
	router := http.NewServeMux()

	handlers := []domain.Handler{
		user.UserHandlers{
			UserService: services.UserService,
		},
	}

	for _, handler := range handlers {
		handler.Init(router)
	}

	return router
}
