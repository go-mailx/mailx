package main

import (
	"context"
	"log"
	"net/mail"

	"github.com/go-mailx/mailx"
	"github.com/go-mailx/mailx-ses"
)

func main() {
	ctx := context.Background()

	// NewFromContext loads AWS config from environment variables, ~/.aws/config, or IAM role.
	factory, err := ses.NewFromContext(ctx)
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
		mailx.Bcc("archive@example.com"),
		mailx.Subject("Hello from go-mailx"),
		mailx.HtmlBody("<h1>Hello!</h1><p>This is a test email sent via Amazon SES.</p>"),
		mailx.TextBody("Hello!\n\nThis is a test email sent via Amazon SES."),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully")
}
