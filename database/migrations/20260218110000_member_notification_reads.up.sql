SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS member_notification_reads (
    member_id uuid REFERENCES members (id) ON DELETE CASCADE,
    notification_id uuid REFERENCES audit_log (id) ON DELETE CASCADE,
    read_at timestamp DEFAULT current_timestamp,
    PRIMARY KEY (member_id, notification_id)
);

--bun:split

CREATE INDEX IF NOT EXISTS member_notification_reads_member_id_idx ON member_notification_reads (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_notification_reads_notification_id_idx ON member_notification_reads (notification_id);
