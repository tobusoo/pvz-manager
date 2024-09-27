-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists refunds_order_id_idx on refunds using btree(order_id);
-- +goose Down
drop index concurrently if exists refunds_order_id_idx;