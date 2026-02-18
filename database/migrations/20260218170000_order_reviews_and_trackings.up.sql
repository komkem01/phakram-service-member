SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS order_shipping_trackings (
    id uuid PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE REFERENCES orders (id) ON DELETE CASCADE,
    tracking_no varchar NOT NULL,
    updated_by uuid REFERENCES members (id),
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS order_shipping_trackings_order_id_idx ON order_shipping_trackings (order_id);

--bun:split

CREATE INDEX IF NOT EXISTS order_shipping_trackings_updated_by_idx ON order_shipping_trackings (updated_by);

--bun:split

CREATE INDEX IF NOT EXISTS order_shipping_trackings_tracking_no_idx ON order_shipping_trackings (tracking_no);

--bun:split

CREATE TABLE IF NOT EXISTS order_payment_reviews (
    id uuid PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE REFERENCES orders (id) ON DELETE CASCADE,
    payment_id uuid REFERENCES payments (id),
    review_status varchar NOT NULL,
    rejected_reason text,
    reviewed_by uuid REFERENCES members (id),
    reviewed_at timestamp,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS order_payment_reviews_order_id_idx ON order_payment_reviews (order_id);

--bun:split

CREATE INDEX IF NOT EXISTS order_payment_reviews_payment_id_idx ON order_payment_reviews (payment_id);

--bun:split

CREATE INDEX IF NOT EXISTS order_payment_reviews_review_status_idx ON order_payment_reviews (review_status);

--bun:split

CREATE INDEX IF NOT EXISTS order_payment_reviews_reviewed_by_idx ON order_payment_reviews (reviewed_by);

--bun:split

WITH latest_tracking AS (
    SELECT DISTINCT ON (al.action_id)
        al.action_id AS order_id,
        TRIM(REPLACE(al.action_detail, 'Shipping tracking number:', '')) AS tracking_no,
        al.action_by,
        al.created_at
    FROM audit_log al
    WHERE al.action_type = 'order_shipping_tracking_updated'
      AND al.status = 'success'
    ORDER BY al.action_id, al.created_at DESC
)
INSERT INTO order_shipping_trackings (id, order_id, tracking_no, updated_by, created_at, updated_at)
SELECT
    uuid_generate_v4(),
    lt.order_id,
    lt.tracking_no,
    lt.action_by,
    lt.created_at,
    lt.created_at
FROM latest_tracking lt
WHERE lt.tracking_no <> ''
ON CONFLICT (order_id) DO UPDATE
SET tracking_no = EXCLUDED.tracking_no,
    updated_by = EXCLUDED.updated_by,
    updated_at = EXCLUDED.updated_at;

--bun:split

WITH latest_payment_review AS (
    SELECT DISTINCT ON (al.action_id)
        al.action_id AS order_id,
        al.action_type,
        al.action_detail,
        al.action_by,
        al.created_at
    FROM audit_log al
    WHERE al.action_type IN ('order_payment_submitted', 'order_payment_approved', 'order_payment_rejected')
      AND al.status = 'success'
    ORDER BY al.action_id, al.created_at DESC
)
INSERT INTO order_payment_reviews (
    id,
    order_id,
    payment_id,
    review_status,
    rejected_reason,
    reviewed_by,
    reviewed_at,
    created_at,
    updated_at
)
SELECT
    uuid_generate_v4(),
    o.id,
    o.payment_id,
    CASE
        WHEN lpr.action_type = 'order_payment_rejected' THEN 'rejected'
        WHEN lpr.action_type = 'order_payment_approved' THEN 'approved'
        ELSE 'submitted'
    END AS review_status,
    CASE
        WHEN lpr.action_type = 'order_payment_rejected' THEN NULLIF(TRIM(REPLACE(lpr.action_detail, 'Payment rejected reason:', '')), '')
        ELSE ''
    END AS rejected_reason,
    CASE
        WHEN lpr.action_type IN ('order_payment_rejected', 'order_payment_approved') THEN lpr.action_by
        ELSE NULL
    END AS reviewed_by,
    CASE
        WHEN lpr.action_type IN ('order_payment_rejected', 'order_payment_approved') THEN lpr.created_at
        ELSE NULL
    END AS reviewed_at,
    lpr.created_at,
    lpr.created_at
FROM latest_payment_review lpr
JOIN orders o ON o.id = lpr.order_id
ON CONFLICT (order_id) DO UPDATE
SET payment_id = EXCLUDED.payment_id,
    review_status = EXCLUDED.review_status,
    rejected_reason = EXCLUDED.rejected_reason,
    reviewed_by = EXCLUDED.reviewed_by,
    reviewed_at = EXCLUDED.reviewed_at,
    updated_at = EXCLUDED.updated_at;
