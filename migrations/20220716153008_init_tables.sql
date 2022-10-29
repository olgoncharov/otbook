-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    user VARCHAR(50) NOT NULL,
    token VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL,
    ttl INT NOT NULL,
    PRIMARY KEY (user),
    FOREIGN KEY (user) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS profiles (
    user VARCHAR(50) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    city VARCHAR(255) NOT NULL,
    sex VARCHAR(10) NOT NULL,
    hobby TEXT NOT NULL,
    PRIMARY KEY (user),
    FOREIGN KEY (user) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS friends (
    user VARCHAR(50) NOT NULL,
    friend VARCHAR(50) NOT NULL,
    FOREIGN KEY (user) REFERENCES users(username),
    FOREIGN KEY (friend) REFERENCES users(username),
    CONSTRAINT unique_friends UNIQUE (user, friend)
);

-- +goose Down
DROP TABLE users;
DROP TABLE refresh_tokens;
DROP TABLE profiles;
DROP TABLE friends;
