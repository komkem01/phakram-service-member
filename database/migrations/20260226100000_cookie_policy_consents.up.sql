SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS cookie_policy_versions (
    id uuid PRIMARY KEY,
    policy_key varchar NOT NULL,
    version_no integer NOT NULL,
    title varchar NOT NULL,
    content text NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    effective_at timestamp NOT NULL DEFAULT current_timestamp,
    created_by uuid REFERENCES members (id),
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp
);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS cookie_policy_versions_policy_key_version_no_uidx
    ON cookie_policy_versions (policy_key, version_no);

--bun:split

CREATE INDEX IF NOT EXISTS cookie_policy_versions_policy_key_is_active_idx
    ON cookie_policy_versions (policy_key, is_active, effective_at DESC);

--bun:split

CREATE TABLE IF NOT EXISTS cookie_policy_consents (
    id uuid PRIMARY KEY,
    policy_version_id uuid NOT NULL REFERENCES cookie_policy_versions (id) ON DELETE CASCADE,
    member_id uuid REFERENCES members (id),
    visitor_key varchar NOT NULL,
    accepted_at timestamp NOT NULL DEFAULT current_timestamp,
    user_agent text,
    created_at timestamp NOT NULL DEFAULT current_timestamp,
    updated_at timestamp NOT NULL DEFAULT current_timestamp
);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS cookie_policy_consents_version_visitor_uidx
    ON cookie_policy_consents (policy_version_id, visitor_key);

--bun:split

CREATE INDEX IF NOT EXISTS cookie_policy_consents_member_id_idx
    ON cookie_policy_consents (member_id);

--bun:split

INSERT INTO cookie_policy_versions (
    id,
    policy_key,
    version_no,
    title,
    content,
    is_active,
    effective_at,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),
    'cookie_notice',
    1,
    'เงื่อนไขการใช้งานคุกกี้',
    'เว็บไซต์นี้ใช้คุกกี้ที่จำเป็นต่อการทำงานของระบบ รวมถึงคุกกี้เพื่อวิเคราะห์การใช้งานและพัฒนาประสบการณ์ของผู้ใช้ การกดยินยอมถือว่าคุณรับทราบและยอมรับเงื่อนไขการใช้งานเว็บไซต์และนโยบายคุกกี้ของเรา',
    true,
    current_timestamp,
    current_timestamp,
    current_timestamp
WHERE NOT EXISTS (
    SELECT 1
    FROM cookie_policy_versions
    WHERE policy_key = 'cookie_notice'
);
