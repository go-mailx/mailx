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

// InlineAttachmentAdder may be implemented by a MailInstance that supports
// inline attachments. Mailer.Send returns ErrInlineAttachmentsNotSupported
// for a mail with inline attachments if the instance does not implement it.
type InlineAttachmentAdder interface {
	AddInlineAttachment(InlineAttachment) error
}
