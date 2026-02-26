SET statement_timeout = 0;

--bun:split

ALTER TABLE storages
ADD COLUMN IF NOT EXISTS file_source varchar;

--bun:split

ALTER TABLE system_bank_accounts
ADD COLUMN IF NOT EXISTS qr_image_source varchar;

--bun:split

UPDATE storages
SET file_source = CASE
    WHEN COALESCE(file_path, '') ILIKE 'data:%' THEN 'INLINE'
    WHEN COALESCE(file_path, '') = '' THEN NULL
    ELSE 'STORAGE'
END
WHERE COALESCE(file_source, '') = '';

--bun:split

UPDATE system_bank_accounts
SET qr_image_source = CASE
    WHEN COALESCE(qr_image_url, '') ILIKE 'data:%' THEN 'INLINE'
    WHEN COALESCE(qr_image_url, '') = '' THEN NULL
    ELSE 'STORAGE'
END
WHERE COALESCE(qr_image_source, '') = '';
