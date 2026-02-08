SET statement_timeout = 0;

--bun:split

DROP TABLE IF EXISTS audit_log;

--bun:split

DROP TYPE IF EXISTS action_type_enum;
DROP TYPE IF EXISTS status_audit_enum;
