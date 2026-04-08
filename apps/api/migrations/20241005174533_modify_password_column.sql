-- +goose Up
-- +goose StatementBegin
ALTER TABLE league_management.users
ALTER COLUMN password TYPE text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
