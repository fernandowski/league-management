-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.teams
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_teams_name ON league_management.teams(name);

CREATE TABLE IF NOT EXISTS league_management.league_teams
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    league_id  UUID NOT NULL REFERENCES league_management.leagues(id),
    team_id    UUID NOT NULL REFERENCES league_management.teams(id),
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_league_teams_league_id_team_id ON league_management.league_teams(league_id, team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS league_management.idx_league_teams_league_id_team_id;
DROP INDEX IF EXISTS league_management.idx_teams_name;
DROP TABLE IF EXISTS league_management.league_teams;
DROP TABLE IF EXISTS league_management.teams;
-- +goose StatementEnd
