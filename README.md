Overview:

This project comprises two service : auth-service and otp-service. 
Together, they handle user authentication and One-Time Password (OTP) management. The services communicate via RPC calls and RabbitMQ.

Technologies:

Programming Language: Go
RPC Framework: Connect RPC
Messaging: RabbitMQ
Database: PG
SMS Management: Twilio API

Architecture:

Auth-Service
Handles user registration, verification, login, and profile retrieval and It also  communicates with the otp-service to generate and verify OTPs.

OTP-Service
Handles OTP generation and sends them to the user using Twilio. It listens to OTP requests and sends OTPs to the auth-service.

RPC Methods:

Auth-Service
.............
1. SignupWithPhoneNumber

Description: Allows users to sign up using their phone number. An OTP is generated and sent to the user's phone for verification.
Req: {
  "phone_number": "+1234567890"
}

Response: {

  "message": "Signup successful, OTP will be sent soon"
}


2. VerifyPhoneNumber

Description: Verifies a user's phone number using the OTP they received and send out the session id.

Req: {
  "phone_number": "+1234567890",
  "otp": "123456"
}

Response: {
  "session_token": "session-token-placeholder"
}


3. LoginWithPhoneNumber

Description :  Logs in users using their phone number and OTP.

Req:  {
  "phone_number": "+1234567890",
  "otp": "123456"
}

Response: {
  "session_token": "session-token-placeholder"
}


4. ValidatePhoneNumberLogin

Description: Checks if a phone number is valid and registered for login.

Req: {
  "phone_number": "+1234567890"
}

Response: {
  "is_valid": true
}

5. GetProfile

Description: Retrieves a user's profile information based on their phone number.

Req: {
  "phone_number": "+1234567890"
}

Response: {
  "phone_number": "+1234567890"
}

OTP-Service
............

1.GenerateOtp

Description: Generates a randomized OTP and sends it to the requesting service, also using Twilio to send the OTP to the user's phone.

Req: {
  "phone_number": "+1234567890"
}

Response: {
  "message": "OTP sent"
}


RabbitMQ Communication
.................


Auth-Service Publisher:

Exchange: verification
Routing Key: otp_routing_key
Description: Sends OTP requests to otp-service.

OTP-Service Consumer:

Queue: otp_queue
Exchange: verification
Routing Key: otp_routing_key
Description: Receives OTP requests and generates OTPs.

OTP-Service Publisher:

Exchange: verification
Routing Key: otp_to_auth
Description: Sends generated OTPs to auth-service.

Auth-Service Consumer:

Queue: auth_service_queue
Exchange: verification
Routing Key: otp_to_auth
Description: Receives OTPs and processes them.


Usage:
.......

Starting the Services
Make sure RabbitMQ is running.
Start the auth-service.
Start the otp-service.
Testing the Services

Use a tool like Postman to send HTTP POST requests to the appropriate endpoints.
Ensure the headers and body match the expected format.