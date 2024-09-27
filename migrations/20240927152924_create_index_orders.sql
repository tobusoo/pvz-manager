-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists orders_user_id_order_id_idx on orders using btree(user_id, order_id);
-- +goose Down
drop index concurrently if exists orders_user_id_order_id_idx;