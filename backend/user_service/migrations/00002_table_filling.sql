-- +goose Up
-- +goose StatementBegin
INSERT INTO users (login, email, password, is_banned) VALUES
('john_doe', 'john.doe@example.com', 'password123', FALSE),
('jane_smith', 'jane.smith@example.com', 'securePass!456', FALSE),
('alex_jones', 'alex.jones@example.com', 'alexPass789', TRUE),
('mary_jane', 'mary.jane@example.com', 'marySecret321', FALSE),
('peter_parker', 'peter.parker@example.com', 'spiderMan!2024', TRUE);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users;
-- +goose StatementEnd
