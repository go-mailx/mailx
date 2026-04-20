package main

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/go-mailx/mailx"
	ses "github.com/go-mailx/mailx-ses"
)

func main() {
	ctx := context.Background()

	// NewFromContext loads AWS config from environment variables, ~/.aws/config, or IAM role.
	adapter, err := ses.NewFromContext(ctx)
	if err != nil {
		slog.Error("failed to create SES adapter", "err", err)
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
		mailx.Bcc("archive@example.com"),
		mailx.Subject("Hello from go-mailx"),
		mailx.HtmlBody("<h1>Hello!</h1><p>This is a test email sent via Amazon SES.</p>"),
		mailx.TextBody("Hello!\n\nThis is a test email sent via Amazon SES."),
	)
	if err != nil {
		slog.Error("failed to send email", "err", err)
		return
	}

	slog.Info("email sent successfully")
}
