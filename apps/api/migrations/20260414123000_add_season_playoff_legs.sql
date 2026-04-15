-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.season_playoff_legs
(
    id            UUID PRIMARY KEY,
    playoff_tie_id UUID NOT NULL REFERENCES league_management.season_playoff_ties (id) ON DELETE CASCADE,
    leg_number    INTEGER NOT NULL,
    home_team_id  UUID REFERENCES league_management.teams (id),
    away_team_id  UUID REFERENCES league_management.teams (id),
    home_score    INTEGER NOT NULL DEFAULT 0,
    away_score    INTEGER NOT NULL DEFAULT 0,
    status        match_status NOT NULL DEFAULT 'scheduled',
    created_at    timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (playoff_tie_id, leg_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.season_playoff_legs;
-- +goose StatementEnd
