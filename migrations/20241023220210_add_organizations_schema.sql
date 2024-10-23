-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.organizations
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(100) UNIQUE NOT NULL,
    user_id    UUID REFERENCES league_management.users (id),
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE league_management.organizations
-- +goose StatementEnd
