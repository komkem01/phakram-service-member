SET statement_timeout = 0;

--bun:split

ALTER TABLE prefixes
	ADD COLUMN IF NOT EXISTS gender_id uuid;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'prefixes_gender_id_fkey'
	) THEN
		ALTER TABLE prefixes
			ADD CONSTRAINT prefixes_gender_id_fkey
			FOREIGN KEY (gender_id) REFERENCES genders(id);
	END IF;
END$$;

--bun:split

CREATE INDEX IF NOT EXISTS prefixes_gender_id_idx ON prefixes (gender_id);
