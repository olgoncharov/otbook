-- +goose Up
CREATE INDEX profile_name_idx ON profiles (last_name, first_name);

-- +goose Down
DROP INDEX profile_name_idx ON profiles;
