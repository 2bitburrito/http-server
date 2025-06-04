-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
   ADD CONSTRAINT user_id UNIQUE ( id );
CREATE TABLE IF NOT EXISTS chirps(
  id UUID NOT NULL DEFAULT uuid_generate_v4(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  body TEXT NOT NULL,
  user_id UUID NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chirps;
-- +goose StatementEnd
