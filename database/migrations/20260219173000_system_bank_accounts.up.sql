CREATE TABLE IF NOT EXISTS system_bank_accounts (
    id uuid PRIMARY KEY,
    bank_id uuid NOT NULL REFERENCES banks(id),
    account_name varchar NOT NULL,
    account_no varchar NOT NULL,
    branch varchar,
    is_active boolean NOT NULL DEFAULT true,
    is_default_receive boolean NOT NULL DEFAULT false,
    is_default_refund boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS system_bank_accounts_account_no_uidx
    ON system_bank_accounts (account_no)
    WHERE account_no IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS system_bank_accounts_default_receive_uidx
    ON system_bank_accounts (is_default_receive)
    WHERE is_default_receive = true;

CREATE UNIQUE INDEX IF NOT EXISTS system_bank_accounts_default_refund_uidx
    ON system_bank_accounts (is_default_refund)
    WHERE is_default_refund = true;

CREATE INDEX IF NOT EXISTS system_bank_accounts_bank_id_idx
    ON system_bank_accounts (bank_id);

CREATE INDEX IF NOT EXISTS system_bank_accounts_active_idx
    ON system_bank_accounts (is_active);
