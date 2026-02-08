SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS districts (
	id uuid PRIMARY KEY,
	province_id uuid REFERENCES provinces(id),
	name varchar,
	is_active bool default true,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS districts_province_id_idx ON districts (province_id);
