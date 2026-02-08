SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS product_files (
	id uuid PRIMARY KEY,
	product_id uuid REFERENCES products(id),
	file_id uuid REFERENCES storages(id),
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS product_files_product_id_idx ON product_files (product_id);
