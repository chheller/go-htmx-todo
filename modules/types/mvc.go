package types

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Init(client *mongo.Client, ctx context.Context) Service
}

type Handler interface {
	Init(router *http.ServeMux) Handler
}
