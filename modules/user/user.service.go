package user

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/chheller/go-htmx-todo/modules/config"
	"github.com/chheller/go-htmx-todo/modules/domain"
	smtp "github.com/chheller/go-htmx-todo/modules/email"
	"github.com/chheller/go-htmx-todo/modules/event"
	"github.com/chheller/go-htmx-todo/modules/web"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService struct {
	client     *mongo.Client
	ctx        context.Context
	collection *mongo.Collection
}

func (svc UserService) Init(client *mongo.Client, ctx context.Context) *UserService {
	svc.client = client
	svc.ctx = ctx
	svc.collection = client.Database("go-todo-htmx").Collection("user")
	return &svc
}

func (svc *UserService) VerifyUserOtp(token string) bool {
	log.Debug("Verifying token")

	mac := hmac.New(sha256.New, []byte("secret"))
	bytesWritten, err := mac.Write([]byte(token))
	if err != nil {
		log.WithField("error", err).Error("Failed to hash token")
		return false
	}
	if bytesWritten != len(token) {
		log.WithField("error", err).Error("Failed to write complete bytes to hash")
		return false
	}

	tokenHash := hex.EncodeToString(mac.Sum(nil))

	query := svc.collection.FindOne(
		svc.ctx,
		bson.M{
			"verificationtoken": tokenHash,
			"$or": bson.A{
				bson.M{
					"event.type": "EmailOtpIssued",
					"expiresat":  bson.M{"$gt": time.Now()},
				},
				bson.M{"event.type": "EmailOtpRevoked"}},
		},
		options.FindOne().SetSort(bson.M{"_id": -1}).SetProjection(bson.D{{"event", 1}, {"userid", 1}}),
	)

	// This only sorta works because the Revoked and Issued events have a common structure
	// TODO: Figure out a better way to handle decoding into a common struct
	var emailOtpEvent EmailOtpIssued
	err = query.Decode(&emailOtpEvent)
	if err != nil {
		log.WithField("error", err).Error("Error fetching email otp event")
		return false

	}

	if emailOtpEvent.Event.Type == "EmailOtpRevoked" {
		log.WithField("event", emailOtpEvent).Debug("Email otp already revoked")
		return false
	}
	err = query.Decode(&emailOtpEvent)
	if err == nil {
		res, err := svc.collection.InsertOne(
			svc.ctx,
			EmailOtpRevoked{
				Event:             event.Event{Timestamp: time.Now(), Type: "EmailOtpRevoked"},
				UserId:            emailOtpEvent.UserId,
				VerificationToken: tokenHash,
			},
		)
		if err != nil {
			log.WithFields(log.Fields{"error": err, "event": emailOtpEvent}).Error("Failed to revoke email otp")
		}
		log.WithField("result", res).Debug("Succesfully revoked email otp")
		return true
	}

	return false
}

var ErrUserInsert = fmt.Errorf("%w:%s", domain.ErrorApplicationGeneric, "failed to insert user")

func (svc *UserService) CreateUser(user User) error {
	log.WithField("user", user).Info("Creating new user")
	userCreatedEvent := &UserCreated{
		Email:  user.Email,
		Event:  event.Event{Timestamp: time.Now(), Type: "UserCreated"},
		UserId: uuid.New(),
	}

	res, err := svc.collection.InsertOne(svc.ctx, userCreatedEvent)

	if err != nil {
		return fmt.Errorf("%w:%s", ErrUserInsert, err)
	}

	// Fire off an email without blocking the request
	// TODO: Error handling- maybe emit an event indicating verification email failed
	go svc.IssueEmailOtp(userCreatedEvent)

	log.WithField("result", res).Info("Successfully created a new user")
	return nil
}

func (svc *UserService) IssueEmailOtp(userCreatedEvent *UserCreated) {
	tokenChallenge, tokenHash, err := createEmailOtp()
	if err != nil {
		log.WithField("error", err).Error("Create verification token error")
		return
	}

	log.WithFields(log.Fields{"tokenHash": tokenHash}).Debug("Created email verification token")

	_, err = svc.collection.InsertOne(svc.ctx, EmailOtpIssued{
		Event:             event.Event{Timestamp: time.Now(), Type: "EmailOtpIssued"},
		UserId:            userCreatedEvent.UserId,
		VerificationToken: tokenHash,
		IssuedAt:          time.Now(),
		ExpiresAt:         time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		log.WithField("error", err).Error("Insert OTP Error")
		return
	}

	redirectUrl := fmt.Sprintf("%s?token=%s", config.GetEnvironment().EmailVerificationRedirectUrl, tokenChallenge)
	var emailBodyBytes bytes.Buffer
	err = web.GetTemplates().ExecuteTemplate(&emailBodyBytes, "email/verify-new-user", VerifyEmailData{RedirectUrl: redirectUrl})
	if err != nil {
		log.WithField("error", err).Error("Template execution error")
		return
	}
	emailBodyString := emailBodyBytes.String()
	if err = smtp.SendEmail(userCreatedEvent.Email, "Verify Email", emailBodyString); err != nil {
		log.WithField("error", err).Error("Send OTP Email error")
		return
	}
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
