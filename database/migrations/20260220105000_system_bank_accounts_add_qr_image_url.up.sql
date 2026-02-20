ALTER TABLE system_bank_accounts
ADD COLUMN IF NOT EXISTS qr_image_url varchar NOT NULL DEFAULT '';
