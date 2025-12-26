DROP INDEX IF EXISTS idx_review_is_approved;

DROP INDEX IF EXISTS idx_review_rating;

DROP INDEX IF EXISTS idx_review_order_id;

DROP INDEX IF EXISTS idx_review_customer_id;

DROP INDEX IF EXISTS idx_review_product_id;

DROP INDEX IF EXISTS idx_variant_is_active;

DROP INDEX IF EXISTS idx_variant_sku;

DROP INDEX IF EXISTS idx_variant_product_id;

DROP INDEX IF EXISTS idx_customer_profile_customer_id;

DROP INDEX IF EXISTS idx_customer_created_at;

DROP INDEX IF EXISTS idx_customer_is_active;

DROP INDEX IF EXISTS idx_customer_uuid;

DROP INDEX IF EXISTS idx_customer_email;

DROP INDEX IF EXISTS idx_category_path;

DROP INDEX IF EXISTS idx_category_level;

DROP INDEX IF EXISTS idx_category_is_active;

DROP INDEX IF EXISTS idx_category_parent_id;

DROP INDEX IF EXISTS idx_category_slug;

DROP INDEX IF EXISTS idx_supplier_is_active;

DROP INDEX IF EXISTS idx_supplier_email;

DROP INDEX IF EXISTS idx_supplier_code;

DROP INDEX IF EXISTS idx_shipping_status;

DROP INDEX IF EXISTS idx_shipping_tracking_number;

DROP INDEX IF EXISTS idx_shipping_order_id;

DROP INDEX IF EXISTS idx_brand_is_active;

DROP INDEX IF EXISTS idx_brand_slug;

DROP INDEX IF EXISTS idx_warehouse_is_active;

DROP INDEX IF EXISTS idx_warehouse_code;

DROP INDEX IF EXISTS idx_order_placed_at;

DROP INDEX IF EXISTS idx_order_shipping_status;

DROP INDEX IF EXISTS idx_order_payment_status;

DROP INDEX IF EXISTS idx_order_status;

DROP INDEX IF EXISTS idx_order_customer_id;

DROP INDEX IF EXISTS idx_order_order_number;

DROP INDEX IF EXISTS idx_product_price;

DROP INDEX IF EXISTS idx_product_brand_id;

DROP INDEX IF EXISTS idx_product_is_featured;

DROP INDEX IF EXISTS idx_product_is_active;

DROP INDEX IF EXISTS idx_product_status;

DROP INDEX IF EXISTS idx_product_slug;

DROP INDEX IF EXISTS idx_product_sku;

DROP INDEX IF EXISTS idx_payment_gateway_transaction_id;

DROP INDEX IF EXISTS idx_payment_status;

DROP INDEX IF EXISTS idx_payment_transaction_id;

DROP INDEX IF EXISTS idx_payment_order_id;

DROP INDEX IF EXISTS idx_order_item_variant_id;

DROP INDEX IF EXISTS idx_order_item_product_id;

DROP INDEX IF EXISTS idx_order_item_order_id;

DROP INDEX IF EXISTS idx_inventory_available_quantity;

DROP INDEX IF EXISTS idx_inventory_warehouse_id;

DROP INDEX IF EXISTS idx_inventory_variant_id;

DROP INDEX IF EXISTS idx_inventory_product_id;

DROP INDEX IF EXISTS idx_address_postal_code;

DROP INDEX IF EXISTS idx_address_country;

DROP INDEX IF EXISTS idx_address_type;

DROP INDEX IF EXISTS idx_address_customer_id;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_order_id;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_customer_id;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_product_id;

ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_product_variants_product_id;

ALTER TABLE customer_profiles DROP CONSTRAINT IF EXISTS fk_customer_profiles_customer_id;

ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_parent_id;

ALTER TABLE shipping DROP CONSTRAINT IF EXISTS fk_shipping_order_id;

ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_shipping_address_id;

ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_billing_address_id;

ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_customer_id;

ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_supplier_id;

ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_brand_id;

ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_order_id;

ALTER TABLE order_items DROP CONSTRAINT IF EXISTS fk_order_items_variant_id;

ALTER TABLE order_items DROP CONSTRAINT IF EXISTS fk_order_items_product_id;

ALTER TABLE order_items DROP CONSTRAINT IF EXISTS fk_order_items_order_id;

ALTER TABLE inventory DROP CONSTRAINT IF EXISTS fk_inventory_warehouse_id;

ALTER TABLE inventory DROP CONSTRAINT IF EXISTS fk_inventory_variant_id;

ALTER TABLE inventory DROP CONSTRAINT IF EXISTS fk_inventory_product_id;

ALTER TABLE addresses DROP CONSTRAINT IF EXISTS fk_addresses_customer_id;

DROP TABLE IF EXISTS reviews CASCADE;

DROP TABLE IF EXISTS product_variants CASCADE;

DROP TABLE IF EXISTS customer_profiles CASCADE;

DROP TABLE IF EXISTS customers CASCADE;

DROP TABLE IF EXISTS categories CASCADE;

DROP TABLE IF EXISTS suppliers CASCADE;

DROP TABLE IF EXISTS shipping CASCADE;

DROP TABLE IF EXISTS brands CASCADE;

DROP TABLE IF EXISTS warehouses CASCADE;

DROP TABLE IF EXISTS orders CASCADE;

DROP TABLE IF EXISTS products CASCADE;

DROP TABLE IF EXISTS payments CASCADE;

DROP TABLE IF EXISTS order_items CASCADE;

DROP TABLE IF EXISTS inventory CASCADE;

DROP TABLE IF EXISTS addresses CASCADE;
