# mailx

[![Go Reference](https://pkg.go.dev/badge/github.com/go-mailx/mailx.svg)](https://pkg.go.dev/github.com/go-mailx/mailx)

`github.com/go-mailx/mailx` is a small Go library for sending email via swappable adapters. It provides a common `MailerAdapter` interface, a functional-options API for composing messages, and a `Mailer` that ties everything together.

## Install

```sh
go get github.com/go-mailx/mailx
```

## Usage

See the [examples](https://github.com/go-mailx/mailx/tree/main/examples) directory for working examples with each adapter.

### Inline images

Images can be embedded in the mail and referenced from the HTML body by their Content-ID. Unlike `data:` URLs, this is rendered by Gmail web and mobile.

```go
err := mailer.Send(ctx,
	mailx.To("user@example.com"),
	mailx.Subject("Welcome"),
	mailx.HtmlBody(`<img src="cid:logo" alt="Logo">`),
	mailx.Inline("logo", "image/png", logoPNG),
)
```

The Content-ID is given without angle brackets and must be unique within the mail. Inline attachments require an adapter whose `MailInstance` implements `mailx.InlineAttachmentAdder`; otherwise `Send` returns `mailx.ErrInlineAttachmentsNotSupported` instead of sending a mail with broken images.

## API

See [`pkg.go.dev/github.com/go-mailx/mailx`](https://pkg.go.dev/github.com/go-mailx/mailx) for all message options, `MailerConfig`, and `FromAddressFunc` built-ins.

## Adapters

| Module                                                                     | Transport                                                         |
| -------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| [`github.com/go-mailx/mailx-smtp`](https://github.com/go-mailx/mailx-smtp) | SMTP via [go-mail](https://github.com/wneessen/go-mail)           |
| [`github.com/go-mailx/mailx-ses`](https://github.com/go-mailx/mailx-ses)   | AWS SES via [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) |
