package mailx

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

type fakeAdapter struct {
	instance MailInstance
}

func (f *fakeAdapter) NewMail(context.Context) (MailInstance, error) { return f.instance, nil }

// basicInstance does not implement InlineAttachmentAdder.
type basicInstance struct {
	sent bool
}

func (b *basicInstance) To([]string) error          { return nil }
func (b *basicInstance) From(string) error          { return nil }
func (b *basicInstance) Bcc([]string) error         { return nil }
func (b *basicInstance) ReplyTo(string) error       { return nil }
func (b *basicInstance) Subject(string) error       { return nil }
func (b *basicInstance) HtmlBody(string) error      { return nil }
func (b *basicInstance) TextBody(string) error      { return nil }
func (b *basicInstance) Send(context.Context) error { b.sent = true; return nil }

type inlineInstance struct {
	basicInstance
	inline []InlineAttachment
}

func (i *inlineInstance) AddInlineAttachment(att InlineAttachment) error {
	i.inline = append(i.inline, att)
	return nil
}

func newTestMailer(instance MailInstance) *Mailer {
	return &Mailer{
		MailerAdapter: &fakeAdapter{instance: instance},
		Config:        &MailerConfig{FromAddressSrc: []FromAddressFunc{StaticFromAddress("from@example.com")}},
	}
}

func TestSendInlineAttachments(t *testing.T) {
	instance := &inlineInstance{}
	err := newTestMailer(instance).Send(context.Background(),
		To("to@example.com"),
		HtmlBody(`<img src="cid:logo"><img src="cid:banner">`),
		Inline("logo", "image/png", []byte("png")),
		Inline("banner", "image/jpeg", []byte("jpeg")),
	)
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if !instance.sent {
		t.Fatal("mail was not sent")
	}
	want := []InlineAttachment{
		{ContentID: "logo", ContentType: "image/png", Data: []byte("png")},
		{ContentID: "banner", ContentType: "image/jpeg", Data: []byte("jpeg")},
	}
	if len(instance.inline) != len(want) {
		t.Fatalf("got %d inline attachments, want %d", len(instance.inline), len(want))
	}
	for i, got := range instance.inline {
		if got.ContentID != want[i].ContentID || got.ContentType != want[i].ContentType || !bytes.Equal(got.Data, want[i].Data) {
			t.Errorf("inline attachment %d = %+v, want %+v", i, got, want[i])
		}
	}
}

func TestSendInlineAttachmentsNotSupported(t *testing.T) {
	instance := &basicInstance{}
	err := newTestMailer(instance).Send(context.Background(),
		To("to@example.com"),
		Inline("logo", "image/png", []byte("png")),
	)
	if !errors.Is(err, ErrInlineAttachmentsNotSupported) {
		t.Fatalf("Send error = %v, want %v", err, ErrInlineAttachmentsNotSupported)
	}
	if instance.sent {
		t.Fatal("mail was sent despite unsupported inline attachments")
	}
}

func TestSendWithoutInlineAttachmentsNeedsNoSupport(t *testing.T) {
	instance := &basicInstance{}
	if err := newTestMailer(instance).Send(context.Background(), To("to@example.com")); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if !instance.sent {
		t.Fatal("mail was not sent")
	}
}

func TestSendInvalidInlineAttachments(t *testing.T) {
	data := []byte("png")
	tests := map[string][]MailOpt{
		"empty content ID":        {Inline("", "image/png", data)},
		"space in content ID":     {Inline("my logo", "image/png", data)},
		"tab in content ID":       {Inline("logo\t", "image/png", data)},
		"newline in content ID":   {Inline("logo\n", "image/png", data)},
		"angle brackets":          {Inline("<logo>", "image/png", data)},
		"less-than in content":    {Inline("lo<go", "image/png", data)},
		"greater-than in content": {Inline("lo>go", "image/png", data)},
		"empty content type":      {Inline("logo", "", data)},
		"nil data":                {Inline("logo", "image/png", nil)},
		"empty data":              {Inline("logo", "image/png", []byte{})},
		"duplicate content ID": {
			Inline("logo", "image/png", data),
			Inline("logo", "image/png", data),
		},
	}
	for name, opts := range tests {
		t.Run(name, func(t *testing.T) {
			instance := &inlineInstance{}
			err := newTestMailer(instance).Send(context.Background(), append([]MailOpt{To("to@example.com")}, opts...)...)
			if !errors.Is(err, ErrInvalidInlineAttachment) {
				t.Fatalf("Send error = %v, want %v", err, ErrInvalidInlineAttachment)
			}
			if len(instance.inline) > 0 || instance.sent {
				t.Fatal("invalid inline attachment reached the adapter")
			}
		})
	}
}
