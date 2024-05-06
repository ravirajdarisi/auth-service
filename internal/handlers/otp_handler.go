package rpc

import (
	"context"
	"fmt"
	"math/rand"
	"ravirajdarisi/auth-service/api/protobufs"
	"ravirajdarisi/auth-service/internal/app/services/otp_service"
	"time"
	"github.com/bufbuild/connect-go"
)

type OtpServiceServer struct{}

func NewOtpServiceServer() *OtpServiceServer {
    return &OtpServiceServer{}
}

func (s *OtpServiceServer) GenerateOtp(ctx context.Context, req *connect.Request[protobufs.OtpRequest]) (*connect.Response[protobufs.OtpResponse], error) {
    phoneNumber := req.Msg.PhoneNumber
    otp := generateRandomOtp()

    if err := otp_service.SendOtpUsingTwilio(phoneNumber, otp); err != nil {
        return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not send OTP: %v", err))
    }

    return connect.NewResponse(&protobufs.OtpResponse{Message: "OTP sent"}), nil
}

func generateRandomOtp() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// func HandleOtpRequest(phoneNumber string) error {
//     otp := generateRandomOtp()
//     err := otp_service.SendOtpUsingTwilio(phoneNumber, otp)
//     if err != nil {
//         return err
//     }
//     // This function will now only handle the sending logic,
//     // while the OTP consumer handles the RabbitMQ communication
//     return nil
// }
