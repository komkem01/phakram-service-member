SET statement_timeout = 0;

--bun:split

ALTER TYPE member_action_enum ADD VALUE IF NOT EXISTS 'read';
