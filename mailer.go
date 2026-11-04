package mailx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	// ErrInlineAttachmentsNotSupported is returned by Mailer.Send for a mail with
	// inline attachments if the adapter's MailInstance does not implement
	// InlineAttachmentAdder.
	ErrInlineAttachmentsNotSupported = errors.New("mailx: adapter does not support inline attachments")
	// ErrInvalidInlineAttachment is wrapped by the errors Mailer.Send returns for
	// an invalid inline attachment.
	ErrInvalidInlineAttachment = errors.New("mailx: invalid inline attachment")
)

type Mailer struct {
	MailerAdapter
	Config *MailerConfig
}

func (m *Mailer) Send(ctx context.Context, opts ...MailOpt) error {
	mail := mailOpts(opts).Create()
	if err := validateInlineAttachments(mail.InlineAttachments); err != nil {
		return err
	}
	a, err := m.NewMail(ctx)
	if err != nil {
		return err
	}
	var inlineAdder InlineAttachmentAdder
	if len(mail.InlineAttachments) > 0 {
		var ok bool
		if inlineAdder, ok = a.(InlineAttachmentAdder); !ok {
			return ErrInlineAttachmentsNotSupported
		}
	}
	if err := a.From(m.Config.GetActualFromAddress(ctx, mail)); err != nil {
		return err
	}
	if err := a.To(mail.To); err != nil {
		return err
	}
	if len(mail.Bcc) > 0 {
		if err := a.Bcc(mail.Bcc); err != nil {
			return err
		}
	}
	if mail.ReplyTo != "" {
		if err := a.ReplyTo(mail.ReplyTo); err != nil {
			return err
		}
	}
	if err := a.Subject(mail.Subject); err != nil {
		return err
	}
	if mail.HtmlBodyFunc != nil {
		html, err := mail.HtmlBodyFunc()
		if err != nil {
			return err
		}
		a.HtmlBody(html)
	}
	if mail.TextBodyFunc != nil {
		text, err := mail.TextBodyFunc()
		if err != nil {
			return err
		}
		a.TextBody(text)
	}
	for _, att := range mail.InlineAttachments {
		if err := inlineAdder.AddInlineAttachment(att); err != nil {
			return err
		}
	}
	return a.Send(ctx)
}

func validateInlineAttachments(attachments []InlineAttachment) error {
	seen := make(map[string]struct{}, len(attachments))
	for i, att := range attachments {
		switch {
		case att.ContentID == "":
			return fmt.Errorf("%w %d: empty content ID", ErrInvalidInlineAttachment, i)
		case strings.ContainsFunc(att.ContentID, func(r rune) bool { return unicode.IsSpace(r) || r == '<' || r == '>' }):
			return fmt.Errorf("%w %q: content ID must not contain whitespace, '<' or '>'", ErrInvalidInlineAttachment, att.ContentID)
		case att.ContentType == "":
			return fmt.Errorf("%w %q: empty content type", ErrInvalidInlineAttachment, att.ContentID)
		case len(att.Data) == 0:
			return fmt.Errorf("%w %q: empty data", ErrInvalidInlineAttachment, att.ContentID)
		}
		if _, ok := seen[att.ContentID]; ok {
			return fmt.Errorf("%w %q: duplicate content ID", ErrInvalidInlineAttachment, att.ContentID)
		}
		seen[att.ContentID] = struct{}{}
	}
	return nil
}
