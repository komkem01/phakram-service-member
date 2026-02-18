SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_type t
        WHERE t.typname = 'status_type_enum'
    ) AND NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_enum e ON t.oid = e.enumtypid
        WHERE t.typname = 'status_type_enum'
          AND e.enumlabel = 'refund_requested'
    ) THEN
        ALTER TYPE status_type_enum ADD VALUE 'refund_requested' AFTER 'paid';
    END IF;
END$$;
