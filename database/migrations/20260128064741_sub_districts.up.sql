SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS sub_districts (
	id uuid PRIMARY KEY,
	district_id uuid REFERENCES districts(id),
	name varchar,
	is_active bool,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS sub_districts_district_id_idx ON sub_districts (district_id);
