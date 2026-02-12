SET statement_timeout = 0;

--bun:split

CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY,
    category_id uuid REFERENCES categories (id),
    name_th varchar,
    name_en varchar,
    product_no varchar,
    price decimal,
    is_active bool DEFAULT false,
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp DEFAULT current_timestamp,
    deleted_at timestamp
);

--bun:split

CREATE INDEX IF NOT EXISTS products_category_id_idx ON products (category_id);

--bun:split

CREATE INDEX IF NOT EXISTS products_is_active_idx ON products (is_active);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS products_product_no_uidx
    ON products (product_no)
    WHERE deleted_at IS NULL;

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS products_category_name_th_uidx
    ON products (category_id, name_th)
    WHERE deleted_at IS NULL;

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS products_category_name_en_uidx
    ON products (category_id, name_en)
    WHERE deleted_at IS NULL;
