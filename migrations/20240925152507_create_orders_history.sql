-- +goose Up
create table if not exists orders_history (
    user_id bigint not null,
    order_id bigint not null,
    expiration_date text not null,
    package_type text not null,
    weight bigint not null,
    cost bigint not null,
    use_tape boolean not null,
    status text not null,
    updated_at text not null,
    primary key(order_id)
);
-- +goose Down
drop table if exists orders_history;