SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS product_details (
    id uuid PRIMARY KEY,
    product_id uuid REFERENCES products (id),
    description text,
    material varchar,
    dimensions varchar,
    weight decimal,
    care_instructions text
);

--bun:split

CREATE INDEX IF NOT EXISTS product_details_product_id_idx ON product_details (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS product_details_product_uidx ON product_details (product_id);
