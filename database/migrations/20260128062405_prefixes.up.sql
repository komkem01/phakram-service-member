SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS prefixes (
	id uuid PRIMARY KEY,
	name_th varchar,
	name_en varchar,
	created_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS prefixes_name_th_idx ON prefixes (name_th);
