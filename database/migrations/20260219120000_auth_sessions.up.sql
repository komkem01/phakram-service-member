SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS auth_sessions (
    id uuid PRIMARY KEY,
    member_id uuid NOT NULL REFERENCES members (id) ON DELETE CASCADE,
    actor_member_id uuid REFERENCES members (id) ON DELETE SET NULL,
    actor_is_admin boolean NOT NULL DEFAULT false,
    is_acting_as boolean NOT NULL DEFAULT false,
    refresh_token_hash text NOT NULL,
    last_activity timestamp NOT NULL DEFAULT current_timestamp,
    refresh_expires_at timestamp NOT NULL,
    revoked_at timestamp,
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS auth_sessions_member_id_idx ON auth_sessions (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS auth_sessions_last_activity_idx ON auth_sessions (last_activity);

--bun:split

CREATE INDEX IF NOT EXISTS auth_sessions_refresh_expires_at_idx ON auth_sessions (refresh_expires_at);

--bun:split

CREATE INDEX IF NOT EXISTS auth_sessions_revoked_at_idx ON auth_sessions (revoked_at);
