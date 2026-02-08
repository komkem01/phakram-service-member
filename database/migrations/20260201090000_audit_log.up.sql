SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'action_audit_enum') THEN
        CREATE TYPE action_audit_enum AS ENUM ('create', 'update', 'delete', 'order', 'payment');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_audit_enum') THEN
        CREATE TYPE status_audit_enum AS ENUM ('success', 'fail');
    END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS audit_log (
    id uuid PRIMARY KEY,
    action action_audit_enum NOT NULL,
    action_type varchar NOT NULL,
    action_id uuid,
    action_by uuid REFERENCES members(id),
    status status_audit_enum NOT NULL,
    action_detail text,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_action_by_idx ON audit_log (action_by);
CREATE INDEX IF NOT EXISTS audit_log_action_type_idx ON audit_log (action_type);
CREATE INDEX IF NOT EXISTS audit_log_action_id_idx ON audit_log (action_id);
