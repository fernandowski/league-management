-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.teams
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_name ON league_management.teams(name);

CREATE TABLE IF NOT EXISTS league_management.league_teams
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    league_id  UUID NOT NULL REFERENCES league_management.leagues(id),
    team_id    UUID NOT NULL REFERENCES league_management.leagues(id),
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX unique_league_id_team_id ON league_management.league_teams(league_id, team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE league_management.teams;
DROP TABLE league_management.league_teams
-- +goose StatementEnd
