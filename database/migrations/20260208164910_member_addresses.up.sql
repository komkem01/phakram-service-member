SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_addresses (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members (id),
	first_name varchar,
	last_name varchar,
	phone varchar,
	is_default bool DEFAULT false,
	address_no varchar,
	village varchar,
	alley varchar,
	sub_district_id uuid REFERENCES sub_districts (id),
	district_id uuid REFERENCES districts (id),
	province_id uuid REFERENCES provinces (id),
	zipcode_id uuid REFERENCES zipcodes (id),
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_addresses_member_id_idx ON member_addresses (member_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_addresses_default_uidx
	ON member_addresses (member_id)
	WHERE is_default IS TRUE AND deleted_at IS NULL;
