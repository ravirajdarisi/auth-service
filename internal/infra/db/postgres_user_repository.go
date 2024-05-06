package db

import (
	"context"
	"fmt"
	"ravirajdarisi/auth-service/internal/domain/entities"

	"github.com/jackc/pgx/v4"
)

type PostgresUserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(conn *pgx.Conn) *PostgresUserRepository {
	return &PostgresUserRepository{conn: conn}
}

//SignupWithPhoneNumber
func (r *PostgresUserRepository) PhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	var exists bool
	err := r.conn.QueryRow(ctx, "SELECT exists (SELECT 1 FROM users WHERE phone_number=$1)", phoneNumber).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check if phone number exists: %v", err)
	}
	return exists, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, phoneNumber string) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO users (phone_number) VALUES ($1)", phoneNumber)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}


func (r *PostgresUserRepository) SaveOTP(phoneNumber, otp string) error {
    query := "UPDATE users SET otp = $1 WHERE phone_number = $2"
    _, err := r.conn.Exec(context.Background(), query, otp, phoneNumber)
    if err != nil {
        return fmt.Errorf("could not save OTP: %v", err)
    }
    return nil
}




//verify phonnumber
func (r *PostgresUserRepository) GetOTPByPhoneNumber(ctx context.Context, phoneNumber string) (string, error) {
    var otp string
    query := "SELECT otp FROM users WHERE phone_number = $1"
    err := r.conn.QueryRow(ctx, query, phoneNumber).Scan(&otp)
    if err != nil {
        if err == pgx.ErrNoRows {
            return "", fmt.Errorf("phone number not found")
        }
        return "", fmt.Errorf("could not retrieve OTP: %v", err)
    }
    return otp, nil
}


func (r *PostgresUserRepository) VerifyUser(ctx context.Context, phoneNumber string) error {
    query := "UPDATE users SET is_verified = TRUE WHERE phone_number = $1"
    _, err := r.conn.Exec(ctx, query, phoneNumber)
    if err != nil {
        return fmt.Errorf("could not verify user: %v", err)
    }
    return nil
}


func (r *PostgresUserRepository) GetUserProfile(ctx context.Context, phoneNumber string) (*entities.UserProfile, error) {
    var profile entities.UserProfile
    query := "SELECT phone_number FROM users WHERE phone_number = $1"
    err := r.conn.QueryRow(ctx, query, phoneNumber).Scan(&profile.PhoneNumber, &profile.Name, &profile.Email)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("could not retrieve profile: %v", err)
    }
    return &profile, nil
}


func (r *PostgresUserRepository) LogUserActivity(ctx context.Context, phoneNumber, eventType string) error {
    query := "INSERT INTO user_activity (phone_number, event_type) VALUES ($1, $2)"
    _, err := r.conn.Exec(ctx, query, phoneNumber, eventType)
    if err != nil {
        return fmt.Errorf("could not log user activity: %v", err)
    }
    return nil
}

