SET statement_timeout = 0;

--bun:split

ALTER TABLE order_items DROP CONSTRAINT IF EXISTS order_items_order_id_fkey;

--bun:split

ALTER TABLE order_items
  ADD CONSTRAINT order_items_order_id_fkey
  FOREIGN KEY (order_id)
  REFERENCES orders (id)
  ON DELETE CASCADE;
