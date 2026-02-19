SET statement_timeout = 0;

--bun:split

ALTER TABLE storages
    ADD COLUMN IF NOT EXISTS deleted_at timestamp;
