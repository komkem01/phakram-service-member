SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS categories (
	id uuid PRIMARY KEY,
	parent_id uuid REFERENCES categories(id),
	name_th varchar,
	name_en varchar,
	is_active bool,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS categories_parent_id_idx ON categories (parent_id);
