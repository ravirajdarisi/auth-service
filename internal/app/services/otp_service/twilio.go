package otp_service

import (
	"fmt"
	"log"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendOtpUsingTwilio sends a given OTP using Twilio's Messaging API
func SendOtpUsingTwilio(phoneNumber, otp string) error {
    // Replace with your Twilio account SID and auth token
    accountSid := "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
    authToken := "f2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

    // Initialize the Twilio client
    client := twilio.NewRestClientWithParams(twilio.ClientParams{
        Username: accountSid,
        Password: authToken,
    })

    params := &twilioApi.CreateMessageParams{}
    params.SetTo(phoneNumber)
    params.SetFrom("YOUR_TWILIO_PHONE_NUMBER") // Replace with your Twilio phone number
    params.SetBody(fmt.Sprintf("Your OTP is: %s", otp))

    // Access the Api service and send the message
    resp, err := client.Api.CreateMessage(params)
    if err != nil {
        return fmt.Errorf("could not send OTP: %v", err)
    }

    // Log the message details for monitoring and troubleshooting
    log.Printf("OTP sent to %s. SID: %s", phoneNumber, *resp.Sid)

    return nil
}
