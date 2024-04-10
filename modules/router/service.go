package router

import (
	"context"

	"github.com/chheller/go-htmx-todo/modules/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationServices struct {
	UserService *user.UserService
}

func (as ApplicationServices) Init(client *mongo.Client, ctx context.Context) ApplicationServices {
	as.UserService = (user.UserService{}).Init(client, context.Background())
	return as
}
