-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance INT NOT NULL DEFAULT 500
);
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_to INT NOT NULL,
    user_from INT,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_to) REFERENCES users(id),
    FOREIGN KEY (user_from) REFERENCES users(id)
);

INSERT INTO users (id, balance)
VALUES
    (1, 999999999),
    (2, 15000),
    (3, 25000),
    (4, 1000),
    (5, 300),
    (6, 3000)
;
INSERT INTO transactions (user_to, user_from, amount)
VALUES
    (1, null, 999999999),
    (2, 1, 15000),
    (3, 1, 25000),
    (4, 1, 1000),
    (5, 1, 300),
    (6, 1, 3000)
;

-- +goose Down
DROP TABLE transactions;
DROP TABLE users;