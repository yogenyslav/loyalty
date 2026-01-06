-- +goose Up
-- +goose StatementBegin
create table loyalty.poll_accrual (
    order_number text primary key references loyalty.order("number") on delete cascade,
    attempt int not null default 0,
    last_polled timestamp not null default current_timestamp,
    next_poll timestamp not null default current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table loyalty.poll_accrual;

-- +goose StatementEnd