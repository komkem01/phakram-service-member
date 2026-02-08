SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS genders (
	id uuid PRIMARY KEY,
	name_th varchar,
	name_en varchar,
	is_active bool default false,
	created_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS genders_name_th_idx ON genders (name_th);
