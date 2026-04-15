-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS league_management.season_playoff_legs
    RENAME TO season_playoff_matches;

ALTER TABLE IF EXISTS league_management.season_playoff_matches
    RENAME COLUMN leg_number TO match_order;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS league_management.season_playoff_matches
    RENAME COLUMN match_order TO leg_number;

ALTER TABLE IF EXISTS league_management.season_playoff_matches
    RENAME TO season_playoff_legs;
-- +goose StatementEnd
