package rpc

import (
	"context"
	"fmt"
	"ravirajdarisi/auth-service/api/protobufs"
	auth_service "ravirajdarisi/auth-service/internal/app/services/auth_service"
	"github.com/bufbuild/connect-go"
)

type AuthServiceServer struct {
	authService *auth_service.AuthService
}

func NewAuthServiceServer(authService *auth_service.AuthService) *AuthServiceServer {
	return &AuthServiceServer{authService: authService}
}

func (s *AuthServiceServer) SignupWithPhoneNumber(ctx context.Context, req *connect.Request[protobufs.SignupRequest]) (*connect.Response[protobufs.SignupResponse], error) {
	if req.Msg.PhoneNumber == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("phone number required"))
	}

    phoneNumber := req.Msg.PhoneNumber

    // Call the service method with the phone number and context
    message, err := s.authService.SignupWithPhoneNumber(ctx, phoneNumber)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    resp := &protobufs.SignupResponse{Message: message}
    return connect.NewResponse(resp), nil
}



func (s *AuthServiceServer) VerifyPhoneNumber(ctx context.Context, req *connect.Request[protobufs.VerifyRequest]) (*connect.Response[protobufs.VerifyResponse], error) {
    phoneNumber := req.Msg.PhoneNumber
    otp := req.Msg.Otp

    // Call the service layer for OTP verification
    message, err := s.authService.VerifyPhoneNumber(ctx, phoneNumber, otp)
    if err != nil {
        return nil, connect.NewError(connect.CodeInvalidArgument, err)
    }

    return connect.NewResponse(&protobufs.VerifyResponse{Message: message}), nil
}



func (s *AuthServiceServer) LoginWithPhoneNumber(ctx context.Context, req *connect.Request[protobufs.LoginRequest]) (*connect.Response[protobufs.LoginResponse], error) {
    phoneNumber := req.Msg.PhoneNumber
    otp := req.Msg.Otp

    // If OTP is empty, just return the stored OTP (initial request)
    if otp == "" {
        storedOTP, err := s.authService.GetStoredOTP(ctx, phoneNumber)
        if err != nil {
            return nil, connect.NewError(connect.CodeInvalidArgument, err)
        }
        return connect.NewResponse(&protobufs.LoginResponse{Message: "OTP sent", Otp: storedOTP}), nil
    }

    // If OTP is provided, verify the login
    sessionToken, err := s.authService.LoginWithPhoneNumber(ctx, phoneNumber, otp)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, err)
    }

    return connect.NewResponse(&protobufs.LoginResponse{Message: "Login successful", SessionToken: sessionToken}), nil
}



func (s *AuthServiceServer) ValidatePhoneNumberLogin(ctx context.Context, req *connect.Request[protobufs.ValidatePhoneNumberRequest]) (*connect.Response[protobufs.ValidatePhoneNumberResponse], error) {
    phoneNumber := req.Msg.PhoneNumber

    // Validate the phone number format
    if len(phoneNumber) < 10 || len(phoneNumber) > 15 {
        return connect.NewResponse(&protobufs.ValidatePhoneNumberResponse{Message: "Invalid phone number format", IsValid: false}), nil
    }

    // Check if the phone number exists in the system
    exists, err := s.authService.PhoneNumberExists(ctx, phoneNumber)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    // Respond based on the existence of the phone number
    if exists {
        return connect.NewResponse(&protobufs.ValidatePhoneNumberResponse{Message: "Phone number is valid", IsValid: true}), nil
    }
    return connect.NewResponse(&protobufs.ValidatePhoneNumberResponse{Message: "Phone number does not exist", IsValid: false}), nil
}


func (s *AuthServiceServer) GetProfile(ctx context.Context, req *connect.Request[protobufs.ProfileRequest]) (*connect.Response[protobufs.ProfileResponse], error) {
    phoneNumber := req.Msg.PhoneNumber

    // Retrieve user profile data from the service layer
    profile, err := s.authService.GetProfile(ctx, phoneNumber)
    if err != nil {
        return nil, connect.NewError(connect.CodeNotFound, err)
    }

    return connect.NewResponse(&protobufs.ProfileResponse{
        PhoneNumber: profile.PhoneNumber,
        Name:        profile.Name,
        Email:       profile.Email,
    }), nil
}

