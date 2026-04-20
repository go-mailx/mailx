package mailx

import "context"

type FromAddressFunc = func(ctx context.Context, mail Mail) string

type MailerConfig struct {
	FromAddressSrc []FromAddressFunc
}

func StaticFromAddress(address string) FromAddressFunc {
	return func(ctx context.Context, mail Mail) string {
		return address
	}
}

func MailOverrideFromAddress() FromAddressFunc {
	return func(ctx context.Context, mail Mail) string {
		if mail.From != nil {
			return mail.From.String()
		}
		return ""
	}
}

func (mc *MailerConfig) GetActualFromAddress(ctx context.Context, mail Mail) string {
	for _, fn := range mc.FromAddressSrc {
		if a := fn(ctx, mail); a != "" {
			return a
		}
	}
	return ""
}
