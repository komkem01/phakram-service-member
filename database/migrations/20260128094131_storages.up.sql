SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'related_entity_enum') THEN
		CREATE TYPE related_entity_enum AS ENUM ('MEMBER_FILE', 'ORDER_FILE', 'PRODUCT_FILE');
	END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS storages (
	id uuid PRIMARY KEY,
	ref_id uuid,
	file_name varchar,
	file_path varchar,
	file_type varchar,
	file_size varchar,
	related_entity related_entity_enum,
	uploaded_by uuid REFERENCES members(id),
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp
);
