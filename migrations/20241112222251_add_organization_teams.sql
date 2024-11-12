-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.organization_teams
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES league_management.organizations (id),
    team_id         UUID NOT NULL REFERENCES league_management.teams (id),
    created_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX unq_idx_organization_id_team_id ON league_management.organization_teams (organization_id, team_id);
CREATE INDEX idx_organization_teams_organization_id ON league_management.organization_teams (organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.organization_teams
-- +goose StatementEnd
