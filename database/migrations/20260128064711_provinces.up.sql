SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS provinces (
	id uuid PRIMARY KEY,
	name varchar,
	is_active bool default true,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS provinces_name_idx ON provinces (name);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS provinces_name_uidx ON provinces (name);
