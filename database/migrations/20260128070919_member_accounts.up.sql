SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_accounts (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members(id),
	email varchar UNIQUE,
	password varchar,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_accounts_member_id_idx ON member_accounts (member_id);
