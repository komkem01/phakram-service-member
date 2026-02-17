SET statement_timeout = 0;

--bun:split

ALTER TABLE cart_items
    ALTER COLUMN price_per_unit DROP NOT NULL,
    ALTER COLUMN total_item_amount DROP NOT NULL;

--bun:split

ALTER TABLE cart_items
    ALTER COLUMN price_per_unit DROP DEFAULT,
    ALTER COLUMN total_item_amount DROP DEFAULT;
