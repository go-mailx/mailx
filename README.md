# mailx

[![Go Reference](https://pkg.go.dev/badge/github.com/go-mailx/mailx.svg)](https://pkg.go.dev/github.com/go-mailx/mailx)

`github.com/go-mailx/mailx` is a small Go library for sending email via swappable adapters. It provides a common `MailerAdapter` interface, a functional-options API for composing messages, and a `Mailer` that ties everything together.

## Install

```sh
go get github.com/go-mailx/mailx
```

## Usage

See the [examples](https://github.com/go-mailx/mailx/tree/main/examples) directory for working examples with each adapter.

## API

See [`pkg.go.dev/github.com/go-mailx/mailx`](https://pkg.go.dev/github.com/go-mailx/mailx) for all message options, `MailerConfig`, and `FromAddressFunc` built-ins.

## Adapters

| Module | Transport |
|---|---|
| [`github.com/go-mailx/mailx-smtp`](../smtp) | SMTP via [go-mail](https://github.com/wneessen/go-mail) |
| [`github.com/go-mailx/mailx-ses`](../ses) | AWS SES via [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) |
