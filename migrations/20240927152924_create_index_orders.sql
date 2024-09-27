-- +goose NO TRANSACTION
-- +goose Up
create index concurrently orders_user_id_order_id_idx on orders using btree(user_id, order_id);
-- +goose Down
drop index concurrently orders_user_id_order_id_idx;