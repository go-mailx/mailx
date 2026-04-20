package main

import (
	"context"
	"log"
	"net/mail"

	"github.com/go-mailx/mailx"
	"github.com/go-mailx/mailx-smtp"
	gomail "github.com/wneessen/go-mail"
)

func main() {
	ctx := context.Background()

	factory, err := smtp.New(smtp.Config{
		Host:      "smtp.example.com",
		Port:      587,
		Username:  "user@example.com",
		Password:  "secret",
		TLSPolicy: gomail.TLSOpportunistic,
	})
	if err != nil {
		log.Fatal(err)
	}

	mailer := mailx.Mailer{
		Factory: factory,
		Config: &mailx.MailerConfig{
			FromAddressSrc: []mailx.FromAddressFunc{
				mailx.MailOverrideFromAddress(),
				mailx.StaticFromAddress("noreply@example.com"),
			},
		},
	}

	err = mailer.Send(ctx,
		mailx.To("alice@example.com", "bob@example.com"),
		mailx.From(mail.Address{Name: "Example App", Address: "app@example.com"}),
		mailx.Subject("Hello from go-mailx"),
		mailx.HtmlBody("<h1>Hello!</h1><p>This is a test email sent via SMTP.</p>"),
		mailx.TextBody("Hello!\n\nThis is a test email sent via SMTP."),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
