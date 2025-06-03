-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS  "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  email TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
