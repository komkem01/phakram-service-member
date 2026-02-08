SET statement_timeout = 0;

--bun:split

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type_enum') THEN
		CREATE TYPE role_type_enum AS ENUM ('customer', 'admin');
	END IF;
END$$;

--bun:split

CREATE TABLE IF NOT EXISTS members (
	id uuid PRIMARY KEY,
	member_no varchar UNIQUE,
	tier_id uuid REFERENCES tiers(id),
	status_id uuid REFERENCES statuses(id),
	prefix_id uuid REFERENCES prefixes(id),
	gender_id uuid REFERENCES genders(id),
	firstname_th varchar,
	lastname_th varchar,
	firstname_en varchar,
	lastname_en varchar,
	role role_type_enum,
	phone varchar UNIQUE,
	total_spent decimal DEFAULT 0,
	current_points int DEFAULT 0,
	registration timestamp,
	last_login timestamp,
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);
