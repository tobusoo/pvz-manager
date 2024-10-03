-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists orders_history_user_order_expiration_idx on orders_history (user_id, order_id) include (expiration_date);
create index concurrently if not exists orders_history_order_status_idx on orders_history (order_id) include (status);
-- +goose Down
drop index concurrently if exists orders_history_user_order_expiration_idx;
drop index concurrently if exists orders_history_order_status_idx;