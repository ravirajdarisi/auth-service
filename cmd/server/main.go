package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"ravirajdarisi/auth-service/api/protobufs/protobufsconnect"
	auth_service "ravirajdarisi/auth-service/internal/app/services/auth_service"
	rpc "ravirajdarisi/auth-service/internal/handlers"
	"ravirajdarisi/auth-service/internal/infra/db"
	"ravirajdarisi/auth-service/internal/infra/rabbitmq"

	"github.com/jackc/pgx/v4"
)

func main() {


	//Initialize RabbitMQ publisher
	publisher, err := rabbitmq.NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	defer publisher.Close()

    // Set up the RabbitMQ exchange and queue
	if err := publisher.SetupRabbitMQ(); err != nil {
		log.Fatalf("Failed to set up RabbitMQ: %v", err)
	}

	
    // Initialize the RabbitMQ Consumer
    consumer, err := rabbitmq.NewOTPConsumer("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
    }
    defer consumer.OTPConsumerClose()

    // Set up the RabbitMQ exchange and queue
    if err := consumer.SetupRabbitMQ(); err != nil {
        log.Fatalf("Failed to set up RabbitMQ: %v", err)
    }

    if err := consumer.OTPConsume("otp_queue"); err != nil {
        panic(fmt.Sprintf("Failed to consume RabbitMQ messages: %v", err))
    }
    
	

    //db related stuff 
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	userRepo := db.NewPostgresUserRepository(conn)



	// Initialize the RabbitMQ Consumer
	authconsumer, err := rabbitmq.NewConsumer("amqp://guest:guest@localhost:5672/", userRepo)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer authconsumer.Close()
	
	// Set up the RabbitMQ exchange and queue
	if err := consumer.SetupRabbitMQ(); err != nil {
		log.Fatalf("Failed to set up RabbitMQ: %v", err)
	}
	
	// Now you can use the consumer to receive messages
	if err := authconsumer.Consume("auth_service_queue"); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}


	authService := auth_service.NewAuthService(userRepo, publisher)
	authServer :=  rpc.NewAuthServiceServer(authService)

	mux := http.NewServeMux()
	mux.Handle(protobufsconnect.NewAuthServiceHandler(authServer))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
