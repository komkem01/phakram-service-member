SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS payment_files (
	id uuid PRIMARY KEY,
	payment_id uuid REFERENCES payments(id),
	file_id uuid REFERENCES storages(id),
	created_at timestamp DEFAULT current_timestamp,
	updated_at timestamp DEFAULT current_timestamp,
	deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS payment_files_payment_id_idx ON payment_files (payment_id);
