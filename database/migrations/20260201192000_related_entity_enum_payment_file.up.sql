SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'related_entity_enum') THEN
		IF NOT EXISTS (
			SELECT 1
			FROM pg_enum
			WHERE enumlabel = 'PAYMENT_FILE'
			  AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'related_entity_enum')
		) THEN
			ALTER TYPE related_entity_enum ADD VALUE 'PAYMENT_FILE';
		END IF;
	END IF;
END$$;
