-- +goose Up
-- +goose StatementBegin
create schema loyalty;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop schema loyalty cascade;

-- +goose StatementEnd