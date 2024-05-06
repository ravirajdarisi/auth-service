package rabbitmq

import (
	"encoding/json"
	"fmt"
	"ravirajdarisi/auth-service/internal/domain/repository"
    "github.com/streadway/amqp"
)

type Consumer struct {
    conn        *amqp.Connection
    channel     *amqp.Channel
    userRepo    repository.UserRepository
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(amqpURL string, userRepo repository.UserRepository) (*Consumer, error) {
    conn, err := amqp.Dial(amqpURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %v", err)
    }

    return &Consumer{
        conn:     conn,
        channel:  ch,
        userRepo: userRepo,
    }, nil
}

// Consume listens to the specified queue for OTP messages
func (c *Consumer) Consume(queueName string) error {
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
            Otp         string `json:"otp"`
        }
        if err := json.Unmarshal(msg.Body, &otpMessage); err != nil {
            fmt.Printf("Error unmarshaling OTP message: %v\n", err)
            continue
        }

        // Save OTP to the database
        if err := c.userRepo.SaveOTP(otpMessage.PhoneNumber, otpMessage.Otp); err != nil {
            fmt.Printf("Error saving OTP: %v\n", err)
        }
    }
    return nil
}


func (c *Consumer) SetupRabbitMQ() error {
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
        "auth_service_queue", // queue name
        true,                 // durable
        false,                // delete when unused
        false,                // exclusive
        false,                // no-wait
        nil,                  // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare queue: %v", err)
    }

    err = c.channel.QueueBind(
        "auth_service_queue", // queue name
        "otp_to_auth",        // routing key
        "verification",       // exchange name
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed to bind queue to exchange: %v", err)
    }

    return nil
}



// Close closes the RabbitMQ connection and channel
func (c *Consumer) Close() {
    c.channel.Close()
    c.conn.Close()
}
