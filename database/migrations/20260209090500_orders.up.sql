SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type_enum') THEN
        CREATE TYPE status_type_enum AS ENUM (
            'pending',
            'paid',
            'shipping',
            'completed',
            'cancelled'
        );
    END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY,
    order_no varchar,
    member_id uuid REFERENCES members (id),
    payment_id uuid REFERENCES payments (id),
    address_id uuid REFERENCES member_addresses (id),
    status status_type_enum,
    total_amount decimal,
    discount_amount decimal,
    net_amount decimal,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS orders_member_id_idx ON orders (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_payment_id_idx ON orders (payment_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_address_id_idx ON orders (address_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_status_idx ON orders (status);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS orders_order_no_uidx ON orders (order_no);
