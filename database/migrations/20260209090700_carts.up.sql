SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS carts (
    id uuid PRIMARY KEY,
    member_id uuid REFERENCES members (id),
    is_active bool DEFAULT false,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS carts_member_id_idx ON carts (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS carts_is_active_idx ON carts (is_active);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS carts_member_active_uidx
    ON carts (member_id)
    WHERE is_active IS TRUE;
