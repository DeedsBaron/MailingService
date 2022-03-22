package store

import (
	"context"
	"net/http"
	"os"

	//"context"
	"github.com/sirupsen/logrus"
	"log"
	"spam/internal/app/config"
	"spam/internal/app/validation"
)

type Storage interface {
	AddUser(ctx context.Context, user *validation.User) error
	UpdateUser(ctx context.Context, user *validation.User, ID int) error
	DeleteUser(ctx context.Context, ID int) error
	AddMailingList(ctx context.Context, ml *validation.MailingList) error
	GeneralStats(ctx context.Context, w http.ResponseWriter) error
	GetDetailedStats(ctx context.Context, w http.ResponseWriter, id string) error
	UpdateMailingList(ctx context.Context, ml *validation.MailingList, ID int) error
	DeleteMailingList(ctx context.Context, ID int) error
	StartSpam(background context.Context, w http.ResponseWriter, m *validation.MailingList, config *config.Config, logger *logrus.Logger)
	GetPhoneNum(ctx context.Context, clientID int, phone *int) error
	SendMessage(spam *Body, messageID int, conf *config.Config, logger *logrus.Logger, ml *validation.MailingList)
}

func InitStorage(config *config.Config) Storage {
	var err error
	var newStorage Storage
	newStorage, err = NewPostgres(config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	logrus.Info("Successfully connected to database")
	return newStorage
}
