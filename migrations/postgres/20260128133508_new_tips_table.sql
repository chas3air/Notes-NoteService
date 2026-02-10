-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes
 (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_private BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE notes REPLICA IDENTITY FULL;

TRUNCATE TABLE notes;

CREATE INDEX idx_notes_user_id ON notes(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
-- +goose StatementEnd
