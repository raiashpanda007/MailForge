package email

import "context"

type EmailServices interface {
	SendMail(ctx context.Context, apikey string, name string, clientEmail string, subject string, body string)
	SendOTP(ctx context.Context, apikey string, name string, clientEmail string, clientID string)
	VerifyOTP(ctx context.Context, apikey string, clientEmail string, clientID string, clientOTPSend string)
}
