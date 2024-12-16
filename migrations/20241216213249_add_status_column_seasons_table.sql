-- +goose Up
-- +goose StatementBegin
CREATE TYPE season_status AS ENUM (
    'pending',
    'planned',
    'in_progress',
    'finished',
    'paused',
    'undefined'
);

ALTER TABLE league_management.seasons
ADD COLUMN status season_status NOT NULL DEFAULT 'undefined';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.seasons DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS season_status;
-- +goose StatementEnd
