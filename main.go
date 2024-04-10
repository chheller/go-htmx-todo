package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/chheller/go-htmx-todo/modules/database"
	"github.com/chheller/go-htmx-todo/modules/router"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	//Create a channel to recieve shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	//Connect to the database, and defer closing the connection.
	client := database.GetMongoClient()
	defer closeDb(client)

	// Create the router, delegating to modules/router/router.go for the implementation.
	services := (router.ApplicationServices{}).Init(client, context.Background())
	httpRouter := router.CreateRouter(&services)
	srv := &http.Server{
		Handler: httpRouter,
		Addr:    fmt.Sprintf(":%v", 8080),
	}
	// Run the server in a never ending goroutine.
	go func() {
		err := srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed")
		} else {
			log.Panicf("Error stopping server %s", err)
		}

	}()
	//Recieve shutdown signals, and try to shutdown gracefully within 10 seconds.
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Printf("Shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error shutting down server %s", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

func closeDb(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("err closing db connection %s", err)
	} else {
		log.Println("db connection gracefully closed")
	}
}

func parseTemplates() (t *template.Template) {
	t = template.Must(template.ParseGlob("./modules/web/templates/*.go.tmpl"))
	t = template.Must(t.ParseGlob("./modules/web/templates/components/*.go.tmpl"))
	t = template.Must(t.ParseGlob("./modules/web/templates/pages/*.go.tmpl"))
	return
}
