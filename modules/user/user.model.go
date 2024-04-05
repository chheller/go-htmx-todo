package user

import (
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/google/uuid"
)

type UserCreated struct {
	Event  event.Event
	Email  string
	UserId uuid.UUID
}

type UserUpdated struct {
	Event               event.Event
	UserId              uuid.UUID
	UserPasswordUpdated struct {
		Password  string
		Salt      string
		Algorithm PasswordAlgorithm
	}
	UserNameUpdated struct {
		Username string
	}
	UserEmailUpdated struct {
		Email string
	}
}

type User struct {
	Email string
}
