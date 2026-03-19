-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tea_blends (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT name_length_check CHECK (LENGTH(name) <= 500),
  CONSTRAINT description_length_check CHECK (LENGTH(description) <= 500)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tea_blends;
-- +goose StatementEnd
