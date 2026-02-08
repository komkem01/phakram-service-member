SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'action_type_enum') THEN
		CREATE TYPE action_type_enum AS ENUM ('created', 'updated', 'deleted', 'logined', 'registered');
	END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS member_transactions (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members(id),
	action action_type_enum,
	details text,
	created_at timestamp DEFAULT current_timestamp
);
