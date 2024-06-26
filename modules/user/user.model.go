package user

import (
	"time"

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

type OtpIssued struct {
	Event             event.Event
	UserId            uuid.UUID
	VerificationToken string
	IssuedAt          time.Time
	ExpiresAt         time.Time
}

type OtpVerified struct {
	Event             event.Event
	UserId            uuid.UUID
	VerificationToken string
}

type OtpRevoked struct {
	Event             event.Event
	UserId            uuid.UUID
	VerificationToken string
}

type VerifyEmailData struct {
	RedirectUrl string
}

const (
	OtpIssuedEvent     = "OtpIssued"
	OtpVerifiedEvent   = "OtpVerified"
	OtpRevokedEvent    = "OtpRevoked"
	UserCreatedEvent   = "UserCreated"
	UserUpdatedEvent   = "UserUpdated"
	UserDeletedEvent   = "UserDeleted"
	UserLoggedInEvent  = "UserLoggedIn"
	UserLoggedOutEvent = "UserLoggedOut"
)
