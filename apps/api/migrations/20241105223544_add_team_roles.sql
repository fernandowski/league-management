-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.league_roles
(
    id              UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES league_management.organizations (id),
    team_id         UUID         NOT NULL REFERENCES league_management.teams (id),
    role            varchar(100) NOT NULL,
    created_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_league_roles_team_id_role ON league_management.league_roles (team_id, role);
CREATE INDEX IF NOT EXISTS idx_league_roles_organization_id_team_id_role ON league_management.league_roles (organization_id, team_id, role);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS league_management.idx_league_roles_team_id_role;
DROP INDEX IF EXISTS league_management.idx_league_roles_organization_id_team_id_role;
DROP TABLE IF EXISTS league_management.league_roles;
-- +goose StatementEnd
