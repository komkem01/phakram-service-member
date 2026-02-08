SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type_enum') THEN
		IF NOT EXISTS (
			SELECT 1
			FROM pg_enum
			WHERE enumlabel = 'returned'
			  AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'status_type_enum')
		) THEN
			ALTER TYPE status_type_enum ADD VALUE 'returned';
		END IF;
	END IF;
END$$;
