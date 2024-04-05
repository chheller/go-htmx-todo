package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chheller/go-htmx-todo/modules/database"
	"github.com/chheller/go-htmx-todo/modules/router"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	//Create a channel to recieve shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	//Connect to the database, and defer closing the connection.
	client := database.GetMongoClient()
	defer CloseDb(client)

	// Create the router, delegating to modules/router/router.go for the implementation.
	httpRouter := router.CreateRouter(client)
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

func CloseDb(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("err closing db connection %s", err)
	} else {
		log.Println("db connection gracefully closed")
	}
}
