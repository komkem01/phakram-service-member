SET statement_timeout = 0;

--bun:split

ALTER TABLE system_bank_accounts
DROP COLUMN IF EXISTS qr_image_source;

--bun:split

ALTER TABLE storages
DROP COLUMN IF EXISTS file_source;
