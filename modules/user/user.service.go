package user

import (
	"context"
	"fmt"
	"time"

	smtp "github.com/chheller/go-htmx-todo/modules/email"
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Client *mongo.Client
	Ctx    context.Context
}

func (svc *UserService) CreateUser(user User) {
	userCollection := svc.Client.Database("go-todo-htmx").Collection("user")
	userCreatedEvent := UserCreated{
		Email:  user.Email,
		Event:  event.Event{Id: 1, Timestamp: time.Now()},
		UserId: uuid.New(),
	}

	res, err := userCollection.InsertOne(svc.Ctx, userCreatedEvent)

	if err != nil {
		panic(err)
	}
	// Fire off an email without blocking the request
	// TODO: Error handling- maybe emit an event indicating verification email failed
	// TODO: Send a login/verification token
	go func() {
		smtp.SendEmail(user.Email, "Welcome to Go Todo", fmt.Sprintf("Welcome to Go Todo, %s", user.Email))
	}()

	fmt.Println("Inserted a single document: ", res.InsertedID)
}
