-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN is_chirpy_red BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN IF EXISTS is_chirpy_red;
-- +goose StatementEnd
