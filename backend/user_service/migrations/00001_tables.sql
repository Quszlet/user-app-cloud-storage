-- +goose Up
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    login VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose Down
DROP TABLE users;