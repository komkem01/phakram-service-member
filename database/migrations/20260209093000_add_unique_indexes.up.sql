SET statement_timeout = 0;

--bun:split

CREATE INDEX IF NOT EXISTS provinces_name_idx ON provinces (name);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS provinces_name_uidx ON provinces (name);

--bun:split

CREATE INDEX IF NOT EXISTS districts_province_id_idx ON districts (province_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS districts_province_name_uidx ON districts (province_id, name);

--bun:split

CREATE INDEX IF NOT EXISTS sub_districts_district_id_idx ON sub_districts (district_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS sub_districts_district_name_uidx ON sub_districts (district_id, name);

--bun:split

CREATE INDEX IF NOT EXISTS zipcodes_sub_districts_id_idx ON zipcodes (sub_districts_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS zipcodes_sub_district_name_uidx ON zipcodes (sub_districts_id, name);

--bun:split

CREATE INDEX IF NOT EXISTS member_addresses_member_id_idx ON member_addresses (member_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_addresses_default_uidx
    ON member_addresses (member_id)
    WHERE is_default IS TRUE AND deleted_at IS NULL;

--bun:split

CREATE INDEX IF NOT EXISTS member_banks_member_id_idx ON member_banks (member_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_banks_member_bank_no_uidx ON member_banks (member_id, bank_no);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_banks_default_uidx
    ON member_banks (member_id)
    WHERE is_default IS TRUE;

--bun:split

CREATE INDEX IF NOT EXISTS member_accounts_member_id_idx ON member_accounts (member_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_accounts_member_id_uidx ON member_accounts (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_wishlist_member_id_idx ON member_wishlist (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_wishlist_product_id_idx ON member_wishlist (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_wishlist_member_product_uidx ON member_wishlist (member_id, product_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_files_member_id_idx ON member_files (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_files_file_id_idx ON member_files (file_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_files_member_file_uidx
    ON member_files (member_id, file_id)
    WHERE deleted_at IS NULL;

--bun:split

CREATE INDEX IF NOT EXISTS member_payments_member_id_idx ON member_payments (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS member_payments_payment_id_idx ON member_payments (payment_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS member_payments_member_payment_uidx ON member_payments (member_id, payment_id);

--bun:split

CREATE INDEX IF NOT EXISTS payments_status_idx ON payments (status);

--bun:split

CREATE INDEX IF NOT EXISTS payments_approved_by_idx ON payments (approved_by);

--bun:split

CREATE INDEX IF NOT EXISTS categories_parent_id_idx ON categories (parent_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS categories_parent_name_th_uidx ON categories (parent_id, name_th);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS categories_parent_name_en_uidx ON categories (parent_id, name_en);

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

--bun:split

CREATE INDEX IF NOT EXISTS product_details_product_id_idx ON product_details (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS product_details_product_uidx ON product_details (product_id);

--bun:split

CREATE INDEX IF NOT EXISTS product_stocks_product_id_idx ON product_stocks (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS product_stocks_product_uidx
    ON product_stocks (product_id)
    WHERE deleted_at IS NULL;

--bun:split

CREATE INDEX IF NOT EXISTS product_files_product_id_idx ON product_files (product_id);

--bun:split

CREATE INDEX IF NOT EXISTS product_files_file_id_idx ON product_files (file_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS product_files_product_file_uidx
    ON product_files (product_id, file_id)
    WHERE deleted_at IS NULL;

--bun:split

CREATE INDEX IF NOT EXISTS orders_member_id_idx ON orders (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_payment_id_idx ON orders (payment_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_address_id_idx ON orders (address_id);

--bun:split

CREATE INDEX IF NOT EXISTS orders_status_idx ON orders (status);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS orders_order_no_uidx ON orders (order_no);

--bun:split

CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);

--bun:split

CREATE INDEX IF NOT EXISTS order_items_product_id_idx ON order_items (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS order_items_order_product_uidx ON order_items (order_id, product_id);

--bun:split

CREATE INDEX IF NOT EXISTS carts_member_id_idx ON carts (member_id);

--bun:split

CREATE INDEX IF NOT EXISTS carts_is_active_idx ON carts (is_active);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS carts_member_active_uidx
    ON carts (member_id)
    WHERE is_active IS TRUE;

--bun:split

CREATE INDEX IF NOT EXISTS cart_items_cart_id_idx ON cart_items (cart_id);

--bun:split

CREATE INDEX IF NOT EXISTS cart_items_product_id_idx ON cart_items (product_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS cart_items_cart_product_uidx ON cart_items (cart_id, product_id);

--bun:split

CREATE INDEX IF NOT EXISTS payment_files_payment_id_idx ON payment_files (payment_id);

--bun:split

CREATE INDEX IF NOT EXISTS payment_files_file_id_idx ON payment_files (file_id);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS payment_files_payment_file_uidx
    ON payment_files (payment_id, file_id)
    WHERE deleted_at IS NULL;

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_action_by_idx ON audit_log (action_by);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_action_type_idx ON audit_log (action_type);

--bun:split

CREATE INDEX IF NOT EXISTS audit_log_status_idx ON audit_log (status);

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS members_member_no_uidx
    ON members (member_no)
    WHERE deleted_at IS NULL;

--bun:split

CREATE UNIQUE INDEX IF NOT EXISTS members_phone_uidx
    ON members (phone)
    WHERE deleted_at IS NULL;
