-- +goose Up
-- +goose StatementBegin
create table loyalty."user"(
    "id" bigserial primary key,
    "login" text not null unique,
    hashed_password text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table loyalty."user";

-- +goose StatementEnd