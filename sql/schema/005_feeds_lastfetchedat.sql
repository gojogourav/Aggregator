-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetches_at TIMESTAMP;

-- +goose Down

ALTER TABLE feeds DROP COLUMN last_fetches_at;