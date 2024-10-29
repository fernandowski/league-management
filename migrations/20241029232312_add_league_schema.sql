-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.leagues
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(100) UNIQUE NOT NULL,
    user_id    UUID REFERENCES league_management.users (id),
    organization_id UUID REFERENCES league_management.organizations(id),
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON league_management.leagues (user_id);
CREATE INDEX idx_organization_id ON league_management.leagues (organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.leagues;
DROP INDEX IF EXISTS  idx_user_id;
DROP INDEX IF EXISTS  idx_organization_id;
-- +goose StatementEnd
