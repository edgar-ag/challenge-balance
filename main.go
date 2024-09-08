package main

import (
	"challenge/balance/models"
	"challenge/balance/notifications"
	"challenge/balance/service"
	"context"
	"errors"
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
		CustomerName:  "Juan Perez",
		Email:         "test@domain.com",
	}
	HandleRequest(context.Background(), event)
}

func HandleRequest(ctx context.Context, event *models.CustomerInfo) {
	config, err := getConfig()
	if err != nil {
		log.Fatalf("Error loading .env file. %v\n", err)
	}
	log.Println("the config was created.")

	txns, err := service.ProccessFile(config.txnsFile)
	if err != nil {
		log.Fatalf("Error trying to proccess the file. %v\n", err)
	}
	log.Println("the file was proccessed.")

	dataTemplate, err := service.GetBalance(txns)
	if err != nil {
		log.Fatalf("Error getting the balance. %v\n", err)
	}
	log.Println("the balance was calculated.")

	dataTemplate.CustomerName = event.CustomerName
	emailClient := notifications.NewEmailClient(config.smtpUrl, config.senderEmail, config.senderEmailPass, dataTemplate)
	err = emailClient.SendNotification(ctx, event)
	if err != nil {
		log.Fatalf("Error sending email. %v\n", err)
	}

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
