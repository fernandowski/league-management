-- +goose Up
-- +goose StatementBegin
CREATE TYPE match_status AS ENUM (
    'scheduled',
    'in_progress',
    'finished',
    'suspended',
    'undefined'
    );

ALTER TABLE league_management.season_schedules ADD COLUMN status match_status NOT NULL DEFAULT 'scheduled';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules DROP COLUMN IF EXISTS match_status;
DROP TYPE IF EXISTS match_status;
-- +goose StatementEnd
