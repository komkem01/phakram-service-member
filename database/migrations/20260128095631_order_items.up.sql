SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS order_items (
	id uuid PRIMARY KEY,
	order_id uuid REFERENCES orders(id),
	product_id uuid REFERENCES products(id),
	quantity int,
	price_per_unit decimal,
	total_item_amount decimal,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);
