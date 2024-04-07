package user

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/chheller/go-htmx-todo/modules/config"
	smtp "github.com/chheller/go-htmx-todo/modules/email"
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Client *mongo.Client
	Ctx    context.Context
}

type VerifyEmailData struct {
	RedirectUrl string
}

func (svc *UserService) CreateUser(user User) error {
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
	go func() {
		// TODO: Make a meaningful token
		redirectUrl := fmt.Sprintf("%s?token=%s", config.GetEnvironment().EmailVerificationRedirectUrl, uuid.New())
		template, err := template.ParseFiles("modules/user/templates/verify_email_template.tmpl")
		// TODO: Error handling- maybe emit an event indicating verification email failed

		if err != nil {
			log.Printf("error parsing email template, %s", err)
			return
		}
		var emailBodyBytes bytes.Buffer
		err = template.ExecuteTemplate(&emailBodyBytes, "verify_email_template.tmpl", VerifyEmailData{RedirectUrl: redirectUrl})
		if err != nil {
			log.Print("error executing email template")
			return
		}
		emailBodyString := emailBodyBytes.String()
		smtp.SendEmail(user.Email, "Verify Email", emailBodyString)
	}()

	fmt.Println("Inserted a single document: ", res.InsertedID)
	return nil
}
