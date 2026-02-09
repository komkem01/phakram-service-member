SET statement_timeout = 0;

--bun:split

ALTER TYPE audit_action_enum ADD VALUE IF NOT EXISTS 'read';
