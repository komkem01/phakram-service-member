SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS order_cancellations (
    id uuid PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE REFERENCES orders (id) ON DELETE CASCADE,
    cancelled_by uuid REFERENCES members (id),
    cancelled_role varchar,
    reason text,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS order_cancellations_order_id_idx ON order_cancellations (order_id);

--bun:split

CREATE INDEX IF NOT EXISTS order_cancellations_cancelled_by_idx ON order_cancellations (cancelled_by);

--bun:split

CREATE INDEX IF NOT EXISTS order_cancellations_cancelled_role_idx ON order_cancellations (cancelled_role);
