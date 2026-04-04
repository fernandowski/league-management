-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.organizations
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.organizations
DROP COLUMN IF EXISTS is_active;
-- +goose StatementEnd
