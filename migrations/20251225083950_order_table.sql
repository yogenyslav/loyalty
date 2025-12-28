-- +goose Up
-- +goose StatementBegin
create type order_status as enum (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED',
    "REGISTERED"
);

create table loyalty.order(
    "number" text primary key,
    "status" order_status not null default 'NEW',
    accrual numeric(10, 2) not null default 0,
    fk_user_id bigserial references loyalty."user"(id) on delete cascade,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table loyalty.order;

drop type order_status;

-- +goose StatementEnd