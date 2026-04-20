package main

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/go-mailx/mailx"
	smtp "github.com/go-mailx/mailx-smtp"
)

func main() {
	ctx := context.Background()

	adapter, err := smtp.New(smtp.Config{
		Host:      "smtp.example.com",
		Port:      587,
		Username:  "user@example.com",
		Password:  "secret",
		TLSPolicy: smtp.TLSOpportunistic,
	})
	if err != nil {
		slog.Error("failed to create SMTP adapter", "err", err)
		return
	}

	mailer := mailx.Mailer{
		MailerAdapter: adapter,
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
		slog.Error("failed to send email", "err", err)
		return
	}

	slog.Info("email sent successfully")
}
