SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS zipcodes (
	id uuid PRIMARY KEY,
	sub_districts_id uuid REFERENCES sub_districts(id),
	name varchar,
	is_active bool default true,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS zipcodes_sub_districts_id_idx ON zipcodes (sub_districts_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS zipcodes_sub_district_name_uidx ON zipcodes (sub_districts_id, name);
