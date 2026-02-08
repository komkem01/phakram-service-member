SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS carts (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members(id),
	is_active bool,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS carts_member_id_idx ON carts (member_id);
