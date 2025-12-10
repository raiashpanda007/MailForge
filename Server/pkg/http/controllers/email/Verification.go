package email

import "context"

type VerificationUtils interface {
	GenerateOTP(ctx context.Context, userID string, apikey string) (string, error)
	VerifyOTP(ctx context.Context, userID string, apikey string, otp string) (bool, error)
}

type OTPRedisDB struct {
}

func NewVerificationUtis() VerificationUtils {
	return OTPRedisDB{}
}
