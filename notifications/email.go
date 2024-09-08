package notifications

import (
	"bytes"
	"challenge/balance/models"
	"context"
	"log"
	"text/template"

	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	SmtpUrl         string
	SenderEmail     string
	SenderEmailPass string
	TemplateData    *models.TemplateEmail
}

func NewEmailClient(smtpUrl, senderEmail, senderEmailPass string, templateData *models.TemplateEmail) *EmailClient {
	return &EmailClient{
		SmtpUrl:         smtpUrl,
		SenderEmail:     senderEmail,
		SenderEmailPass: senderEmailPass,
		TemplateData:    templateData,
	}
}

func (e *EmailClient) SendNotification(ctx context.Context, customerInfo *models.CustomerInfo) error {
	body, err := createTemplate(e.TemplateData)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", e.SenderEmail)
	mail.SetHeader("To", customerInfo.Email)
	mail.SetHeader("Subject", "Balance Info")
	mail.SetBody("text/html", body.String())

	dial := gomail.NewDialer(e.SmtpUrl, 587, e.SenderEmail, e.SenderEmailPass)
	err = dial.DialAndSend(mail)
	if err != nil {
		return err
	}
	log.Println("email was sent successfully")
	return nil
}

func createTemplate(templateData *models.TemplateEmail) (bytes.Buffer, error) {
	var body bytes.Buffer
	templateData.ImageSrc = "./notifications/stori.png"
	emailTemplate, err := template.ParseFiles("./notifications/template.html")
	if err != nil {
		return body, err
	}

	err = emailTemplate.Execute(&body, templateData)
	if err != nil {
		return body, err
	}
	return body, nil
}
