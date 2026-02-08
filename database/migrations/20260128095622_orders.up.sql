SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type_enum') THEN
		CREATE TYPE status_type_enum AS ENUM ('pending', 'paid', 'shipping', 'completed', 'cancelled');
	END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS orders (
	id uuid PRIMARY KEY,
	order_no varchar,
	member_id uuid REFERENCES members(id),
	payment_id uuid REFERENCES payments(id),
	address_id uuid REFERENCES member_addresses(id),
	status status_type_enum,
	total_amount decimal,
	discount_amount decimal,
	net_amount decimal,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);
