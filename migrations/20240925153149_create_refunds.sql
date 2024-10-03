-- +goose Up
create table if not exists refunds (
    order_id bigint not null,
    primary key(order_id)
);
-- +goose Down
drop table if exists refunds;