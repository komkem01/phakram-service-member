SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_audit_enum') THEN
        CREATE TYPE status_audit_enum AS ENUM (
            'successes',
            'failed'
        );
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'audit_action_enum') THEN
        CREATE TYPE audit_action_enum AS ENUM (
            'created',
            'updated',
            'deleted',
            'logined',
            'registered'
        );
    END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS audit_log (
    id uuid PRIMARY KEY,
    action audit_action_enum,
    action_type varchar,
    action_id uuid,
    action_by uuid REFERENCES members (id),
    status status_audit_enum,
    action_detail text,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_action_by_idx ON audit_log (action_by);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_action_type_idx ON audit_log (action_type);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_status_idx ON audit_log (status);
