-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'match_status'
          AND n.nspname = 'public'
    ) THEN
        CREATE TYPE match_status AS ENUM (
            'scheduled',
            'in_progress',
            'finished',
            'suspended',
            'undefined'
        );
    END IF;
END $$;

ALTER TABLE league_management.season_schedules ADD COLUMN IF NOT EXISTS status match_status NOT NULL DEFAULT 'scheduled';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.season_schedules DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS match_status;
-- +goose StatementEnd
