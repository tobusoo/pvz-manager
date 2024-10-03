-- +goose Up
create table if not exists orders (
    user_id bigint not null,
    order_id bigint not null,
    primary key(order_id)
);
-- +goose Down
drop table if exists orders;