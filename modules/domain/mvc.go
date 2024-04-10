package domain

import (
	"net/http"
)

type Handler interface {
	Init(router *http.ServeMux) Handler
}
