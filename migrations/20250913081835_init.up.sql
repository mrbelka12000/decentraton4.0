-- +goose Up
-- +goose StatementBegin

CREATE TABLE clients (
    client_code INT PRIMARY KEY,
    name VARCHAR(512),
    status VARCHAR(512),
    age INT,
    city VARCHAR(512),
    avg_monthly_balance_KZT INT
);

CREATE TABLE transactions (
    client_code	INT REFERENCES clients(client_code),
    name	VARCHAR(512),
    product	VARCHAR(512),
    status	VARCHAR(512),
    city	VARCHAR(512),
    date	VARCHAR(512),
    category	VARCHAR(512),
    amount	DECIMAL,
    currency	VARCHAR(512)
);

CREATE TABLE transfers (
    client_code	INT REFERENCES clients(client_code),
    name	VARCHAR(512),
    product	VARCHAR(512),
    status	VARCHAR(512),
    city	VARCHAR(512),
    date	VARCHAR(512),
    type	VARCHAR(512),
    direction	VARCHAR(512),
    amount	DECIMAL,
    currency	VARCHAR(512)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE transactions;

DROP TABLE clients;
-- +goose StatementEnd
