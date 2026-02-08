SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS banks (
	id uuid PRIMARY KEY,
	name_th varchar,
	name_abb_th varchar,
	name_en varchar,
	name_abb_en varchar,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS banks_name_th_idx ON banks (name_th);
