-- +goose Up
ALTER TABLE profiles
    ADD COLUMN is_celebrity BOOL NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE profiles
    DROP COLUMN is_celebrity;
