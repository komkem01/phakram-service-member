SET statement_timeout = 0;

--bun:split

DROP INDEX IF EXISTS product_files_product_file_uidx;

--bun:split

DROP INDEX IF EXISTS product_files_file_id_idx;

--bun:split

DROP INDEX IF EXISTS product_files_product_id_idx;

--bun:split

DROP TABLE IF EXISTS product_files;
