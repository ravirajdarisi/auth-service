package rabbitmq

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"ravirajdarisi/auth-service/internal/app/services/otp_service"
	"time"

	"github.com/streadway/amqp"
)

type OTPConsumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

// NewConsumer creates a new RabbitMQ consumer
func NewOTPConsumer(amqpURL string) (*OTPConsumer, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %v", err)
    }

    return &OTPConsumer{
        conn:    conn,
        channel: ch,
    }, nil
}

// Consume listens to the specified queue for OTP messages
func (c *OTPConsumer) OTPConsume(queueName string) error {
    msgs, err := c.channel.Consume(
        queueName, // queue
        "",        // consumer
        true,      // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %v", err)
    }

    for msg := range msgs {
        var otpMessage struct {
            PhoneNumber string `json:"phone_number"`
        }
        if err := json.Unmarshal(msg.Body, &otpMessage); err != nil {
            fmt.Printf("Error unmarshaling OTP message: %v\n", err)
            continue
        }

        otp := generateRandomOtp()
        

		// Send OTP using Twilio
        if err := otp_service.SendOtpUsingTwilio(otpMessage.PhoneNumber, otp); err != nil {
            fmt.Printf("Error sending OTP via Twilio: %v\n", err)
            continue
        }
    
        // Publish the OTP to the auth-service
        if err := c.publishOtpToAuthService(otpMessage.PhoneNumber, otp); err != nil {
            fmt.Printf("Error publishing OTP to auth-service: %v\n", err)
        }
    }
    return nil
}


// SetupRabbitMQ configures the exchange and queue for the consumer
func (c *OTPConsumer) SetupRabbitMQ() error {
    // Declare the exchange
    err := c.channel.ExchangeDeclare(
        "verification", // exchange name
        "direct",       // type of exchange
        true,           // durable
        false,          // auto-delete
        false,          // internal
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    // Declare the queue and bind it to the exchange with the routing key
    _, err = c.channel.QueueDeclare(
        "otp_queue", // queue name
        true,        // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %v", err)
    }

    err = c.channel.QueueBind(
        "otp_queue",       // queue name
        "otp_routing_key", // routing key
        "verification",    // exchange name
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue to exchange: %v", err)
    }

    return nil
}


func generateRandomOtp() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}


// Publish OTP to auth-service
func (c *OTPConsumer) publishOtpToAuthService(phoneNumber, otp string) error {
    // Construct OTP message
    otpMessage := struct {
        PhoneNumber string `json:"phone_number"`
        Otp         string `json:"otp"`
    }{
        PhoneNumber: phoneNumber,
        Otp:         otp,
    }

    // Convert the message to JSON
    body, err := json.Marshal(otpMessage)
    if err != nil {
        return fmt.Errorf("could not marshal OTP message: %v", err)
    }

    // Publish the message to the "verification" exchange with the appropriate routing key
    err = c.channel.Publish(
        "verification",           // exchange
        "otp_to_auth",            // routing key
        false,                    // mandatory
        false,                    // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        body,
        })
    if err != nil {
        return fmt.Errorf("could not publish OTP message: %v", err)
    }
    return nil
}


// Close closes the RabbitMQ connection and channel
func (c *OTPConsumer) OTPConsumerClose() {
    c.channel.Close()
    c.conn.Close()
}
