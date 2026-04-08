-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA league_management;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA league_management;
-- +goose StatementEnd
