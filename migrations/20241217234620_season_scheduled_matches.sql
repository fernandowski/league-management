-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS league_management.season_schedules
(
    id              UUID PRIMARY KEY,
    season_id       UUID    NOT NULL REFERENCES league_management.seasons (id),
    league_id       UUID    NOT NULL REFERENCES league_management.leagues (id),
    round           INTEGER NOT NULL,
    home_team_id    UUID    NOT NULL REFERENCES league_management.teams (id),
    away_team_id    UUID    NOT NULL REFERENCES league_management.teams (id),
    home_team_score INTEGER NOT NULL,
    away_team_score INTEGER NOT NULL,
    created_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_season_schedules_season_id ON league_management.season_schedules(season_id);
CREATE INDEX idx_season_schedules_league_id ON league_management.season_schedules(league_id);
CREATE INDEX idx_season_schedules_home_team_id ON league_management.season_schedules(home_team_id);
CREATE INDEX idx_season_schedules_away_team_id ON league_management.season_schedules(away_team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS league_management.season_schedules;
-- +goose StatementEnd
