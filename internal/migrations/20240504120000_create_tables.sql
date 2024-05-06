-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone_number TEXT UNIQUE NOT NULL,
    otp TEXT,
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE user_activity (
    id SERIAL PRIMARY KEY,
    phone_number TEXT NOT NULL,
    event_type TEXT NOT NULL, -- 'login' or 'logout'
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (phone_number) REFERENCES users(phone_number) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS user_activity;
DROP TABLE IF EXISTS users;
