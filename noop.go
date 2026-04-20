package mailx

import "context"

type Noop struct{}

func (n *Noop) NewMail(context.Context) (MailInstance, error) { return n, nil }

func (n *Noop) Bcc([]string) error         { return nil }
func (n *Noop) From(string) error          { return nil }
func (n *Noop) HtmlBody(string) error      { return nil }
func (n *Noop) ReplyTo(string) error       { return nil }
func (n *Noop) Send(context.Context) error { return nil }
func (n *Noop) Subject(string) error       { return nil }
func (n *Noop) TextBody(string) error      { return nil }
func (n *Noop) To([]string) error          { return nil }

var _ MailerAdapter = (*Noop)(nil)
