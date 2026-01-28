-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tips (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_private BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_tips_user_id ON tips(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tips;
-- +goose StatementEnd
