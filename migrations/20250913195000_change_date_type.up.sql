-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
ALTER COLUMN date TYPE DATE USING date::DATE;

ALTER TABLE transfers
ALTER COLUMN date TYPE DATE USING date::DATE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
ALTER COLUMN date TYPE VARCHAR(512) USING date::VARCHAR;

ALTER TABLE transfers
ALTER COLUMN date TYPE VARCHAR(512) USING date::VARCHAR;
-- +goose StatementEnd
