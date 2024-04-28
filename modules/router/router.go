package router

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chheller/go-htmx-todo/modules/config"
	"github.com/chheller/go-htmx-todo/modules/domain"
	"github.com/chheller/go-htmx-todo/modules/user"
	"github.com/chheller/go-htmx-todo/modules/web"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
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
	router.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	if config.GetEnvironment().InjectBrowserReload {
		// Add no-op websocket handler for the reloadBrowser.js functionality
		router.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
			log.Info("Websocket connection established")
			defer log.Info("Websocket connection closed")
			defer ws.Close()

			// Create a channel to listen on the websocket connection
			websocketClosed := make(chan uint)

			// Start a goroutine to read messages from the websocket connection
			go func() {
				for {
					err := websocket.Message.Receive(ws, &struct{}{})
					if err != nil {
						log.Info("Websocket connection closed by client")
						websocketClosed <- 0
						break
					}
				}
			}()

			// Wait for the websocket connection to close
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)
			signal.Notify(stop, syscall.SIGTERM)
			select {
			case <-websocketClosed:
				{

				}
			case <-stop:
				{

				}
			}
		}))
	}

	// Setup application handlers
	handlers := []domain.Handler{
		user.UserHandlers{
			UserService: services.UserService,
		},
		&web.ErrorPageHandlers{},
	}

	for _, handler := range handlers {
		handler.Init(router)
	}

	return router
}
