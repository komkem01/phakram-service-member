SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_type_enum') THEN
		CREATE TYPE payment_type_enum AS ENUM ('pending', 'success', 'failed', 'refunded');
	END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS payments (
	id uuid PRIMARY KEY,
	amount decimal,
	status payment_type_enum,
	approved_by uuid REFERENCES members(id),
	approved_at timestamp
);
