-- +goose Up
-- +goose StatementBegin
create table loyalty.balance(
    fk_user_id bigserial primary key references loyalty."user"("id") on delete cascade,
    "current" numeric(10, 2) not null default 0,
    withdrawn numeric(10, 2) not null default 0,
    updated_at timestamp not null default current_timestamp
);

create table loyalty.withdrawal(
    "id" bigserial primary key,
    amount numeric(10, 2) not null,
    fk_order_number bigint references loyalty.order("number") on delete cascade,
    fk_user_id bigserial references loyalty."user"("id") on delete cascade,
    created_at timestamp not null default current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table loyalty.withdrawal;

drop table loyalty.balance;

-- +goose StatementEnd