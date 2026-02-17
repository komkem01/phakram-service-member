SET statement_timeout = 0;

--bun:split

UPDATE cart_items
SET
    price_per_unit = COALESCE(price_per_unit, 0),
    total_item_amount = COALESCE(total_item_amount, 0)
WHERE price_per_unit IS NULL OR total_item_amount IS NULL;

--bun:split

ALTER TABLE cart_items
    ALTER COLUMN price_per_unit SET DEFAULT 0,
    ALTER COLUMN total_item_amount SET DEFAULT 0;

--bun:split

ALTER TABLE cart_items
    ALTER COLUMN price_per_unit SET NOT NULL,
    ALTER COLUMN total_item_amount SET NOT NULL;
