-- +goose NO TRANSACTION
-- +goose Up
create index concurrently orders_history_order_id_idx on orders_history using btree(order_id);
create index concurrently orders_history_user_id_idx on orders_history using btree(user_id);
-- +goose Down
drop index concurrently orders_history_user_id_idx;
drop index concurrently orders_history_order_id_idx;