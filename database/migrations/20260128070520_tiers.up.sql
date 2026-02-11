SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS tiers (
	id uuid PRIMARY KEY,
	name_th varchar,
	name_en varchar,
	min_spending decimal,
	discount_rate decimal,
	is_active bool default true,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS tiers_name_th_idx ON tiers (name_th);
