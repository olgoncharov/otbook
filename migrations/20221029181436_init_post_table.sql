-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
    id BIGINT NOT NULL AUTO_INCREMENT,
    author VARCHAR(50) NOT NULL,
    title VARCHAR(250) NOT NULL,
    text TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id),
    FOREIGN KEY (author) REFERENCES users(username)
);

-- +goose Down
DROP TABLE posts;
