package user

import (
	"context"
	"fmt"
	"time"

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
		email:  user.email,
		Event:  event.Event{Id: 1, Timestamp: time.Now()},
		UserId: uuid.New(),
	}
	res, err := userCollection.InsertOne(svc.Ctx, userCreatedEvent)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted a single document: ", res.InsertedID)
}
