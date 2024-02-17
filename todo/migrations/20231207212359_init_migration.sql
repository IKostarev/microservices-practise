-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS todo (
    id UUID PRIMARY KEY,
    created_by int NOT NULL,
    assignee int NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todo;
-- +goose StatementEnd
