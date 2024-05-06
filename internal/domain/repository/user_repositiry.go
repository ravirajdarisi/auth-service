package repository

import (
	"context"
	"ravirajdarisi/auth-service/internal/domain/entities"
)

type UserRepository interface {
	PhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error)
	CreateUser(ctx context.Context, phoneNumber string) error
	GetOTPByPhoneNumber(ctx context.Context, phoneNumber string) (string, error)
    VerifyUser(ctx context.Context, phoneNumber string) error
	SaveOTP(phoneNumber, otp string) error
	GetUserProfile(ctx context.Context, phoneNumber string) (*entities.UserProfile, error)
	LogUserActivity(ctx context.Context, phoneNumber, eventType string) error
}
