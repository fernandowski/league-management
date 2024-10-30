-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.organizations
ADD COLUMN is_active BOOLEAN DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE league_management.organizations
DROP COLUMN is_active;
-- +goose StatementEnd
