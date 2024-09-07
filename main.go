package main

import (
	"challenge/balance/service"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Event struct {
	AccountNumber string
}

type config struct {
	dbName   string
	dbHost   string
	dbPort   string
	dbUser   string
	dbPass   string
	TxnsFile string
}

func main() {
	event := &Event{AccountNumber: "123456789012345678"}
	HandleRequest(context.Background(), event)
}

func HandleRequest(ctx context.Context, event *Event) {
	config, err := getConfig()
	if err != nil {
		log.Fatalf("Error loading .env file. %v\n", err)
	}
	log.Println("the config was created.")

	txns, err := service.ProccessFile(config.TxnsFile)
	if err != nil {
		log.Fatalf("Error trying to proccess the file. %v\n", err)
	}
	log.Println("the file was proccessed.")

	data, err := service.GetBalance(txns)
	if err != nil {
		log.Fatalf("Error getting the balance. %v\n", err)
	}
	log.Println("the balance was calculated.")

	fmt.Println(data)
}

func getConfig() (*config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	database := os.Getenv("DATABASE_NAME")
	if database == "" {
		return nil, errors.New("data base name is required")
	}
	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		return nil, errors.New("data base host is required")
	}
	port := os.Getenv("DATABASE_PORT")
	if port == "" {
		return nil, errors.New("port is required")
	}
	user := os.Getenv("DATABASE_USER")
	if user == "" {
		return nil, errors.New("user is required")
	}
	password := os.Getenv("DATABASE_PASSWORD")
	if password == "" {
		return nil, errors.New("password is required")
	}
	txnsFile := os.Getenv("TXNS_FILE")
	if txnsFile == "" {
		return nil, errors.New("txns file is required")
	}

	return &config{
		dbName:   database,
		dbHost:   host,
		dbPort:   port,
		dbUser:   user,
		dbPass:   password,
		TxnsFile: txnsFile,
	}, nil
}
