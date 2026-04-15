-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.season_playoff_brackets
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season_id  UUID NOT NULL UNIQUE REFERENCES league_management.seasons (id) ON DELETE CASCADE,
    status     VARCHAR(32) NOT NULL DEFAULT 'generated',
    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS league_management.season_playoff_qualified_teams
(
    playoff_bracket_id UUID NOT NULL REFERENCES league_management.season_playoff_brackets (id) ON DELETE CASCADE,
    team_id            UUID NOT NULL REFERENCES league_management.teams (id),
    seed               INTEGER NOT NULL,
    PRIMARY KEY (playoff_bracket_id, team_id),
    UNIQUE (playoff_bracket_id, seed)
);

CREATE TABLE IF NOT EXISTS league_management.season_playoff_ties
(
    id                UUID PRIMARY KEY,
    playoff_bracket_id UUID NOT NULL REFERENCES league_management.season_playoff_brackets (id) ON DELETE CASCADE,
    round_name        VARCHAR(64) NOT NULL,
    round_order       INTEGER NOT NULL,
    slot_order        INTEGER NOT NULL,
    home_seed         INTEGER,
    away_seed         INTEGER,
    home_team_id      UUID REFERENCES league_management.teams (id),
    away_team_id      UUID REFERENCES league_management.teams (id),
    status            VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at        timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (playoff_bracket_id, round_order, slot_order)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.season_playoff_ties;
DROP TABLE IF EXISTS league_management.season_playoff_qualified_teams;
DROP TABLE IF EXISTS league_management.season_playoff_brackets;
-- +goose StatementEnd
