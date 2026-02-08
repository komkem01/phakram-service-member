SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS products (
	id uuid PRIMARY KEY,
	category_id uuid REFERENCES categories(id),
	name_th varchar,
	name_en varchar,
	product_no varchar UNIQUE,
	price decimal,
	is_active bool,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS products_category_id_idx ON products (category_id);
