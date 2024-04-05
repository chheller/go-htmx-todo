package main

import (
	"context"
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

type User struct {
	id         int64
	created_at time.Time
	email      string
	user_name  string
	password   string
}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	client := database.GetMongoClient()
	defer CloseDb(client)

	httpRouter := router.CreateRouter(client)
	srv := &http.Server{
		Handler: httpRouter,
		Addr:    fmt.Sprintf(":%v", 8080),
	}
	go func() {
		panic(srv.ListenAndServe())
	}()
	//Recieve shutdown signals.
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
