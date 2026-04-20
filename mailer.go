package mailx

import (
	"context"
)

type Mailer struct {
	Factory MailerAdapter
	Config  *MailerConfig
}

func (m *Mailer) Send(ctx context.Context, opts ...MailOpt) error {
	mail := mailOpts(opts).Create()
	a, err := m.Factory.NewMail(ctx)
	if err != nil {
		return err
	}
	if err := a.From(m.Config.GetActualFromAddress(ctx, mail)); err != nil {
		return err
	}
	if err := a.To(mail.To); err != nil {
		return err
	}
	if err := a.Bcc(mail.Bcc); err != nil {
		return err
	}
	if err := a.ReplyTo(mail.ReplyTo); err != nil {
		return err
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
	return a.Send(ctx)
}
