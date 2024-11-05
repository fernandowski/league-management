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

CREATE UNIQUE INDEX unq_id_team_id_role ON league_management.league_roles (team_id, role);
CREATE INDEX idx_organization_id_team_id_role ON league_management.league_roles (organization_id, team_id, role);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.league_roles;
-- +goose StatementEnd
