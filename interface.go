package mailx

import (
	"context"
)

type MailerAdapter interface {
	NewMail(context.Context) (MailInstance, error)
}

type MailInstance interface {
	To([]string) error
	From(string) error
	Bcc([]string) error
	ReplyTo(string) error
	Subject(string) error
	HtmlBody(string) error
	TextBody(string) error
	Send(context.Context) error
}
