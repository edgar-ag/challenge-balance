package main

import (
	"challenge/balance/database"
	"challenge/balance/models"
	"challenge/balance/notifications"
	"challenge/balance/repository"
	"challenge/balance/service"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	dbName          string
	dbHost          string
	dbPort          string
	dbUser          string
	dbPass          string
	txnsFile        string
	smtpUrl         string
	senderEmail     string
	senderEmailPass string
}

func main() {
	event := &models.CustomerInfo{
		AccountNumber: "123456789012345678",
		Name:          "Juan Perez",
		Email:         "test@domain.com",
	}
	HandleRequest(context.Background(), event)
}

func HandleRequest(ctx context.Context, event *models.CustomerInfo) {
	config, err := getConfig()
	if err != nil {
		log.Fatalf("error loading .env file. %v\n", err)
	}

	serviceBalance := service.NewBalance(config.txnsFile, event)
	txns, err := serviceBalance.ProccessFile()
	if err != nil {
		log.Fatalf("error trying to proccess the file. %v\n", err)
	}
	balanceInfo, err := serviceBalance.GetBalanceInfo(txns)
	if err != nil {
		log.Fatalf("error getting balance info. %v\n", err)
	}

	repo, err := createRepository(config)
	if err != nil {
		log.Fatalf("error creating conection with DB. %v\n", err)
	}
	repository.SetRepository(repo)

	emailClientConfig := &notifications.Config{
		SmtpUrl:         config.smtpUrl,
		SenderEmail:     config.senderEmail,
		SenderEmailPass: config.senderEmailPass,
	}
	emailClient := notifications.NewEmailClient(emailClientConfig, balanceInfo)

	go serviceBalance.InsertDataIntoDB(ctx, txns)

	go emailClient.SendNotification(ctx, event)
}

func createRepository(config *config) (*database.MysqlRepository, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.dbUser, config.dbPass, config.dbHost, config.dbPort, config.dbName)
	repository, err := database.NewMysqlRepository(dbUrl)
	if err != nil {
		return nil, err
	}
	return repository, nil
}

func getConfig() (*config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		return nil, errors.New("data base name is required")
	}
	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost == "" {
		return nil, errors.New("data base host is required")
	}
	dbPort := os.Getenv("DATABASE_PORT")
	if dbPort == "" {
		return nil, errors.New("data base port is required")
	}
	dbUser := os.Getenv("DATABASE_USER")
	if dbUser == "" {
		return nil, errors.New("data base user is required")
	}
	dbPass := os.Getenv("DATABASE_PASSWORD")
	if dbPass == "" {
		return nil, errors.New("data base password is required")
	}
	txnsFile := os.Getenv("TXNS_FILE")
	if txnsFile == "" {
		return nil, errors.New("txns file is required")
	}
	smtpUrl := os.Getenv("SMTP_URL")
	if smtpUrl == "" {
		return nil, errors.New("smtp url is required")
	}
	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		return nil, errors.New("sender email is required")
	}
	senderEmailPass := os.Getenv("SENDER_EMAIL_PASS")
	if senderEmailPass == "" {
		return nil, errors.New("sender email password file is required")
	}

	log.Println("the config was created.")
	return &config{
		dbName:          dbName,
		dbHost:          dbHost,
		dbPort:          dbPort,
		dbUser:          dbUser,
		dbPass:          dbPass,
		txnsFile:        txnsFile,
		smtpUrl:         smtpUrl,
		senderEmail:     senderEmail,
		senderEmailPass: senderEmailPass,
	}, nil
}
