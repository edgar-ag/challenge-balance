package notifications

import (
	"bytes"
	"challenge/balance/models"
	"context"
	"log"
	"sync"
	"text/template"

	"gopkg.in/gomail.v2"
)

type Config struct {
	SmtpUrl         string
	SenderEmail     string
	SenderEmailPass string
}

type EmailClient struct {
	Config      *Config
	BalanceInfo *models.BalanceInfo
}

func NewEmailClient(config *Config, balanceInfo *models.BalanceInfo) *EmailClient {
	return &EmailClient{
		Config:      config,
		BalanceInfo: balanceInfo,
	}
}

func (e *EmailClient) SendNotification(ctx context.Context, wg *sync.WaitGroup, customerInfo *models.CustomerInfo) {
	defer wg.Done()

	body, err := e.createTemplate()
	if err != nil {
		log.Printf("error sending email. %v\n", err)
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", e.Config.SenderEmail)
	mail.SetHeader("To", customerInfo.Email)
	mail.SetHeader("Subject", "Balance Info")
	mail.SetBody("text/html", body.String())

	dial := gomail.NewDialer(e.Config.SmtpUrl, 587, e.Config.SenderEmail, e.Config.SenderEmailPass)
	err = dial.DialAndSend(mail)
	if err != nil {
		log.Printf("error sending email. %v\n", err)
	} else {
		log.Println("email was sent successfully")
	}
}

func (e *EmailClient) createTemplate() (bytes.Buffer, error) {
	var body bytes.Buffer
	e.BalanceInfo.ImageSrc = "./notifications/stori.png"
	emailTemplate, err := template.ParseFiles("./notifications/template.html")
	if err != nil {
		return body, err
	}

	err = emailTemplate.Execute(&body, e.BalanceInfo)
	if err != nil {
		return body, err
	}

	log.Println("the email template was created successfully")
	return body, nil
}
