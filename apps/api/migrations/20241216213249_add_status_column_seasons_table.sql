-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'season_status'
          AND n.nspname = 'public'
    ) THEN
        CREATE TYPE season_status AS ENUM (
            'pending',
            'planned',
            'in_progress',
            'finished',
            'paused',
            'undefined'
        );
    END IF;
END $$;

ALTER TABLE league_management.seasons
ADD COLUMN IF NOT EXISTS status season_status NOT NULL DEFAULT 'undefined';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.seasons DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS season_status;
-- +goose StatementEnd
