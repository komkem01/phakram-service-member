SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_banks (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members (id),
	bank_id uuid REFERENCES banks (id),
	bank_no varchar,
	firstname_th varchar,
	lastname_th varchar,
	firstname_en varchar,
	lastname_en varchar,
	is_default bool DEFAULT false,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_banks_member_id_idx ON member_banks (member_id);
