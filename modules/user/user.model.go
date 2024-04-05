package user

import (
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/google/uuid"
)

type UserCreated struct {
	Event  event.Event
	email  string
	UserId uuid.UUID
}

type UserUpdated struct {
	event               event.Event
	userId              uuid.UUID
	UserPasswordUpdated struct {
		password  string
		salt      string
		algorithm PasswordAlgorithm
	}
	UserNameUpdated struct {
		username string
	}
	UserEmailUpdated struct {
		email string
	}
}

type User struct {
	email string
}
