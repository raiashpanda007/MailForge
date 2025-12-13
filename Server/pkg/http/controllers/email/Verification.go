package email

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/big"
	"time"
)

type OTPStatus int

const (
	OTPNotExists OTPStatus = iota
	OTPBlocked
	OTPAllowed
	OTPError
)

const WindowDuration = 90 * time.Second

type VerificationUtils interface {
	GenerateOTP(ctx context.Context, userID string, apikey string) (*string, error)
	VerifyOTP(ctx context.Context, userID string, apikey string, otp string) (bool, error)
}

type OTPRedisDB struct {
	redisClient *redis.Client
}

type OTPInfo struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
}

func NewVerificationUtis(redisClient *redis.Client) VerificationUtils {
	return &OTPRedisDB{redisClient: redisClient}
}

func (c *OTPRedisDB) OTPStatusCheck(ctx context.Context, userID string, apikey string) (OTPStatus, *string, error) {
	key := fmt.Sprintf("OTP/%s/%s", apikey, userID)
	value, err := c.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return OTPNotExists, nil, nil
	}

	if err != nil {
		return OTPError, nil, err
	}

	jsonData := []byte(value)
	isValid := json.Valid(jsonData)
	if !isValid {
		return OTPError, nil, errors.New("INVALID JSON DATA FROM REDIS FOR OTP VERIFICATION RELATED USERID :: " + userID)
	}
	var otpData OTPInfo
	err = json.Unmarshal(jsonData, &otpData)
	if err != nil {
		return OTPError, nil, err
	}
	windowTime := time.Since(otpData.CreatedAt)

	if windowTime < WindowDuration {
		return OTPBlocked, nil, errors.New("HANG TIGHT! YOU CAN REQUEST A NEW OTP IN A MOMENT.")
	}

	return OTPAllowed, &otpData.Value, nil
}

func (c *OTPRedisDB) GenerateOTP(ctx context.Context, userID string, apikey string) (*string, error) {
	otpStatus, existingOTP, err := c.OTPStatusCheck(ctx, userID, apikey)
	switch otpStatus {
	case OTPBlocked:
		return nil, err
	case OTPAllowed:
		return existingOTP, nil
	case OTPError:
		return nil, err
	}

	var otp string
	n, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return nil, err
	}
	otp = n.String()
	otpInfo := OTPInfo{
		CreatedAt: time.Now(),
		Value:     otp,
	}
	keyString := fmt.Sprintf("OTP/%s/%s", apikey, userID)
	otpJson, err := json.Marshal(otpInfo)
	if err != nil {
		return nil, err
	}
	err = c.redisClient.Set(ctx, keyString, otpJson, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (c *OTPRedisDB) VerifyOTP(ctx context.Context, userID string, apikey string, otp string) (bool, error) {

	return true, nil
}
