SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS statuses (
	id uuid PRIMARY KEY,
	name_th varchar,
	name_en varchar,
	is_active bool default true,
	created_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS statuses_name_th_idx ON statuses (name_th);
