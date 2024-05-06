package auth_service

import (
	"context"
	"encoding/json"
	"fmt"
	"ravirajdarisi/auth-service/internal/domain/entities"
	"ravirajdarisi/auth-service/internal/domain/repository"
	"ravirajdarisi/auth-service/internal/infra/rabbitmq"
)

type AuthService struct {
	userRepository repository.UserRepository
	publisher      *rabbitmq.Publisher
}

func NewAuthService(userRepo repository.UserRepository, publisher *rabbitmq.Publisher ) *AuthService {
	return &AuthService{userRepository: userRepo, publisher: publisher}
}

func (s *AuthService) SignupWithPhoneNumber(ctx context.Context, phoneNumber string) (string, error) {
    // Check if the phone number already exists
    exists, err := s.userRepository.PhoneNumberExists(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not check phone number: %v", err)
    }
    if exists {
        return "", fmt.Errorf("phone number already exists")
    }

    // Create and save the new user
    err = s.userRepository.CreateUser(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not create user: %v", err)
    }

    // Publish SendOTP message to RabbitMQ
    otpMessage := map[string]string{"phone_number": phoneNumber}
    body, err := json.Marshal(otpMessage)
    if err != nil {
        return "", fmt.Errorf("could not marshal OTP message: %v", err)
    }

    // Publish OTP message to RabbitMQ
    err = s.publisher.Publish("verification", "otp_routing_key", body)
    if err != nil {
        return "", fmt.Errorf("could not publish OTP message: %v", err)
    }

    return "Signup successful, OTP will be sent soon", nil
}


func (s *AuthService) VerifyPhoneNumber(ctx context.Context, phoneNumber, otp string) (string, error) {
    // Retrieve the stored OTP from the database
    storedOTP, err := s.userRepository.GetOTPByPhoneNumber(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not retrieve OTP: %v", err)
    }

    // Compare the stored OTP with the provided OTP
    if storedOTP != otp {
        return "", fmt.Errorf("invalid OTP")
    }

    // Mark the user as verified 
    err = s.userRepository.VerifyUser(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not verify user: %v", err)
    }

    sessionToken := "session-token-placeholder"
    return sessionToken, nil
}


func (s *AuthService) GetStoredOTP(ctx context.Context, phoneNumber string) (string, error) {
    otp, err := s.userRepository.GetOTPByPhoneNumber(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not get OTP: %v", err)
    }
    return otp, nil
}

func (s *AuthService) LoginWithPhoneNumber(ctx context.Context, phoneNumber, otp string) (string, error) {
    // Retrieve and verify the OTP
    storedOTP, err := s.userRepository.GetOTPByPhoneNumber(ctx, phoneNumber)
    if err != nil {
        return "", fmt.Errorf("could not get OTP: %v", err)
    }
    if storedOTP != otp {
        return "", fmt.Errorf("invalid OTP")
    }

    // Log the login event
    err = s.userRepository.LogUserActivity(ctx, phoneNumber, "login")
    if err != nil {
        return "", fmt.Errorf("could not log login activity: %v", err)
    }

    // Create a session token (for simplicity, just returning a placeholder)
    sessionToken := "session-token-placeholder"
    return sessionToken, nil
}



func (s *AuthService) PhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
    return s.userRepository.PhoneNumberExists(ctx, phoneNumber)
}


func (s *AuthService) GetProfile(ctx context.Context, phoneNumber string) (*entities.UserProfile, error) {
    return s.userRepository.GetUserProfile(ctx, phoneNumber)
}


