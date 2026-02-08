SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_files (
	id uuid PRIMARY KEY,
	member_id uuid REFERENCES members(id),
	file_id uuid REFERENCES storages(id),
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS member_files_member_id_idx ON member_files (member_id);
