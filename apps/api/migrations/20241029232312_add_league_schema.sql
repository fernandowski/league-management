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

CREATE INDEX IF NOT EXISTS idx_leagues_user_id ON league_management.leagues (user_id);
CREATE INDEX IF NOT EXISTS idx_leagues_organization_id ON league_management.leagues (organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS league_management.idx_leagues_user_id;
DROP INDEX IF EXISTS league_management.idx_leagues_organization_id;
DROP TABLE IF EXISTS league_management.leagues;
-- +goose StatementEnd
