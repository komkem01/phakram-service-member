SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_payments (
    id uuid PRIMARY KEY,
    member_id uuid REFERENCES members (id),
    payment_id uuid REFERENCES payments (id),
    quantity int,
    price decimal,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_payments_member_id_idx ON member_payments (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_payments_payment_id_idx ON member_payments (payment_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_payments_member_payment_uidx ON member_payments (member_id, payment_id);
