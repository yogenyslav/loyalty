-- +goose Up
-- +goose StatementBegin
create schema loyalty;

alter database dev
set
    time zone 'Europe/Moscow';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop schema loyalty cascade;

alter database dev
set
    time zone 'UTC';

-- +goose StatementEnd