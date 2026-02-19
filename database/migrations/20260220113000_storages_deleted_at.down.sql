SET statement_timeout = 0;

--bun:split

ALTER TABLE storages
    DROP COLUMN IF EXISTS deleted_at;
