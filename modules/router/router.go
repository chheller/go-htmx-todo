package router

import (
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/domain"
	"github.com/chheller/go-htmx-todo/modules/user"
	"github.com/chheller/go-htmx-todo/modules/web"
	log "github.com/sirupsen/logrus"
)

// CreateRouter creates a new http.ServeMux and initializes the http handlers.
// Wires up services with a database client, and attaches those services to handlers
// Responsible for middleware stack as well.
func CreateRouter(services *ApplicationServices) *http.ServeMux {
	router := http.NewServeMux()

	// Serve files from static folder
	static, err := web.GetStaticWebAssets()
	if err != nil {
		log.WithError(err).Panic("Failed to create sub filesytem for static files")
	}

	staticFileServer := http.FileServer(http.FS(static))
	log.Println("Serving static files")
	router.Handle("/static/", staticFileServer)

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
