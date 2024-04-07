package user

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"time"

	"github.com/chheller/go-htmx-todo/modules/config"
	smtp "github.com/chheller/go-htmx-todo/modules/email"
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Client *mongo.Client
	Ctx    context.Context
}

type VerifyEmailData struct {
	RedirectUrl string
}

func (svc *UserService) VerifyUserOtp(token string) (ok bool) {
	userCollection := svc.Client.Database("go-todo-htmx").Collection("user")
	mac := hmac.New(sha256.New, []byte("secret"))
	bytesWritten, err := mac.Write([]byte(token))
	if err != nil {
		return
	}
	if bytesWritten != len(token) {
		err = fmt.Errorf("error hashing token")
		return
	}
	tokenHash := hex.EncodeToString(mac.Sum(nil))

	var emailOtpIssued EmailOtpIssued
	var emailOtpRevoked EmailOtpRevoked

	err = userCollection.FindOne(svc.Ctx, bson.M{"verificationtoken": tokenHash, "event.type": "EmailOtpIssued", "expiresat": bson.M{"$gt": time.Now()}}).Decode(&emailOtpIssued)
	if err != nil {
		// TODO: Check if the error is no docs found, return a 404
		log.Print("Did not find any matching Verification Token")
		return
	}
	// TODO: Handle better token already used error
	err = userCollection.FindOne(svc.Ctx, bson.M{"verificationtoken": tokenHash, "event.type": "EmailOtpRevoked"}).Decode(&emailOtpRevoked)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ok = true
			err = nil // nil out err, as not finding a revoked event is actually a success state for verification
			userCollection.InsertOne(svc.Ctx, EmailOtpRevoked{Event: event.Event{Timestamp: time.Now(), Type: "EmailOtpRevoked"}, UserId: emailOtpIssued.UserId, VerificationToken: tokenHash})
			return
		}
		return
	}
	return
}
func (svc *UserService) CreateUser(user User) error {
	userCollection := svc.Client.Database("go-todo-htmx").Collection("user")
	userCreatedEvent := UserCreated{
		Email:  user.Email,
		Event:  event.Event{Timestamp: time.Now(), Type: "UserCreated"},
		UserId: uuid.New(),
	}

	res, err := userCollection.InsertOne(svc.Ctx, userCreatedEvent)

	if err != nil {
		panic(err)
	}

	// Fire off an email without blocking the request
	// TODO: Error handling- maybe emit an event indicating verification email failed
	go func() {
		tokenChallenge, tokenHash, err := createEmailOtp()
		if err != nil {
			log.Printf("error creating email verification token %s", err)
			return
		}

		_, err = userCollection.InsertOne(svc.Ctx, EmailOtpIssued{
			Event:             event.Event{Timestamp: time.Now(), Type: "EmailOtpIssued"},
			UserId:            userCreatedEvent.UserId,
			VerificationToken: tokenHash,
			IssuedAt:          time.Now(),
			ExpiresAt:         time.Now().Add(time.Hour * 24),
		})
		if err != nil {
			log.Printf("error inserting email otp issued %s", err)
			return
		}

		redirectUrl := fmt.Sprintf("%s?token=%s", config.GetEnvironment().EmailVerificationRedirectUrl, tokenChallenge)
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

var maxLetterSize = big.NewInt(51)

// TODO: Consider making this more efficient via masking https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, maxLetterSize)
		if err != nil {
			// TODO: Handle this properly (how can this even fail?)
			panic("Failed to generated random number")
		}
		b[i] = byte(num.Uint64() + 65) // 65 is the ASCII code for 'A', so we want to start there
	}
	return string(b)
}
func createEmailOtp() (tokenChallenge string, tokenHash string, err error) {

	tokenSize := 32
	// TODO: Make this part of the service New function
	mac := hmac.New(sha256.New, []byte("secret"))

	tokenChallenge = randStringBytes(tokenSize)
	bytesHashed, err := mac.Write([]byte(tokenChallenge))
	if err != nil {
		log.Printf("error hashing random bytes %s", err)
		return
	}
	if bytesHashed != tokenSize {
		log.Printf("error hashing random bytes, expected %d, got %d", tokenSize, bytesHashed)
		err = fmt.Errorf("error hashing random bytes")
		return
	}

	tokenHash = hex.EncodeToString(mac.Sum(nil))
	return
}
