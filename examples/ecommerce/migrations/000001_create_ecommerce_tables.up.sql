-- Create table: addresses
CREATE TABLE IF NOT EXISTS addresses (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "customer_id" BIGINT NOT NULL,
    "type" VARCHAR(20) NOT NULL DEFAULT 'shipping',
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "company" VARCHAR(200),
    "address_line1" VARCHAR(255) NOT NULL,
    "address_line2" VARCHAR(255),
    "city" VARCHAR(100) NOT NULL,
    "state" VARCHAR(100) NOT NULL,
    "postal_code" VARCHAR(20) NOT NULL,
    "country" VARCHAR(2) NOT NULL,
    "phone" VARCHAR(20),
    "is_default" BOOLEAN DEFAULT FALSE,
    "is_active" BOOLEAN DEFAULT TRUE,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: inventory
CREATE TABLE IF NOT EXISTS inventory (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "product_id" BIGINT NOT NULL,
    "variant_id" BIGINT,
    "warehouse_id" BIGINT NOT NULL,
    "quantity" INTEGER DEFAULT 0,
    "reserved_quantity" INTEGER DEFAULT 0,
    "available_quantity" INTEGER DEFAULT 0,
    "cost_per_unit" NUMERIC(12, 2),
    "location" VARCHAR(100),
    "is_active" BOOLEAN DEFAULT TRUE,
    "last_restocked_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: order_items
CREATE TABLE IF NOT EXISTS order_items (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" BIGINT NOT NULL,
    "product_id" BIGINT NOT NULL,
    "variant_id" BIGINT,
    "product_name" VARCHAR(500) NOT NULL,
    "product_sku" VARCHAR(100) NOT NULL,
    "quantity" INTEGER NOT NULL,
    "unit_price" NUMERIC(12, 2) NOT NULL,
    "total_price" NUMERIC(12, 2) NOT NULL,
    "discount_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "tax_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "product_data" JSONB,
    "variant_data" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: payments
CREATE TABLE IF NOT EXISTS payments (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" BIGINT NOT NULL,
    "transaction_id" UUID NOT NULL UNIQUE,
    "payment_method" VARCHAR(50) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "amount" NUMERIC(12, 2) NOT NULL,
    "currency" VARCHAR(3) DEFAULT 'USD',
    "gateway" VARCHAR(100),
    "gateway_transaction_id" VARCHAR(255),
    "gateway_response" VARCHAR(50),
    "gateway_message" TEXT,
    "gateway_data" JSONB,
    "card_last4" VARCHAR(4),
    "card_brand" VARCHAR(50),
    "card_expiry" DATE,
    "is_refunded" BOOLEAN DEFAULT FALSE,
    "refund_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "paid_at" TIMESTAMP WITH TIME ZONE,
    "refunded_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: products
CREATE TABLE IF NOT EXISTS products (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "sku" UUID NOT NULL UNIQUE,
    "name" VARCHAR(500) NOT NULL,
    "slug" VARCHAR(500) NOT NULL UNIQUE,
    "description" TEXT,
    "short_description" TEXT,
    "price" NUMERIC(12, 2) NOT NULL,
    "compare_at_price" NUMERIC(12, 2),
    "cost_price" NUMERIC(12, 2),
    "currency" VARCHAR(3) DEFAULT 'USD',
    "stock_quantity" INTEGER DEFAULT 0,
    "low_stock_threshold" INTEGER DEFAULT 10,
    "track_inventory" BOOLEAN DEFAULT TRUE,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_featured" BOOLEAN DEFAULT FALSE,
    "is_digital" BOOLEAN DEFAULT FALSE,
    "requires_shipping" BOOLEAN DEFAULT TRUE,
    "weight" DOUBLE PRECISION,
    "length" DOUBLE PRECISION,
    "width" DOUBLE PRECISION,
    "height" DOUBLE PRECISION,
    "status" VARCHAR(50) DEFAULT 'draft',
    "attributes" JSONB,
    "images" JSONB,
    "seo_data" JSONB,
    "brand_id" BIGINT,
    "supplier_id" BIGINT,
    "published_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: orders
CREATE TABLE IF NOT EXISTS orders (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_number" UUID NOT NULL UNIQUE,
    "customer_id" BIGINT NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'pending',
    "subtotal" NUMERIC(12, 2) NOT NULL,
    "tax_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "shipping_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "discount_amount" NUMERIC(12, 2) DEFAULT 0.00,
    "total_amount" NUMERIC(12, 2) NOT NULL,
    "currency" VARCHAR(3) DEFAULT 'USD',
    "billing_address_id" BIGINT NOT NULL,
    "shipping_address_id" BIGINT NOT NULL,
    "payment_status" VARCHAR(50) DEFAULT 'pending',
    "shipping_status" VARCHAR(50) DEFAULT 'pending',
    "notes" TEXT,
    "customer_notes" TEXT,
    "metadata" JSONB,
    "placed_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "confirmed_at" TIMESTAMP WITH TIME ZONE,
    "shipped_at" TIMESTAMP WITH TIME ZONE,
    "delivered_at" TIMESTAMP WITH TIME ZONE,
    "cancelled_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: warehouses
CREATE TABLE IF NOT EXISTS warehouses (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(200) NOT NULL,
    "code" VARCHAR(50) NOT NULL UNIQUE,
    "address_line1" VARCHAR(255) NOT NULL,
    "address_line2" VARCHAR(255),
    "city" VARCHAR(100) NOT NULL,
    "state" VARCHAR(100) NOT NULL,
    "postal_code" VARCHAR(20) NOT NULL,
    "country" VARCHAR(2) NOT NULL,
    "phone" VARCHAR(20),
    "email" VARCHAR(255),
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_primary" BOOLEAN DEFAULT FALSE,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: brands
CREATE TABLE IF NOT EXISTS brands (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(200) NOT NULL UNIQUE,
    "slug" VARCHAR(250) NOT NULL UNIQUE,
    "description" TEXT,
    "logo_url" VARCHAR(500),
    "website_url" VARCHAR(500),
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_featured" BOOLEAN DEFAULT FALSE,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: shipping
CREATE TABLE IF NOT EXISTS shipping (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" BIGINT NOT NULL UNIQUE,
    "carrier" VARCHAR(100),
    "service" VARCHAR(100),
    "tracking_number" VARCHAR(100),
    "tracking_url" VARCHAR(500),
    "status" VARCHAR(50) DEFAULT 'pending',
    "cost" NUMERIC(12, 2) DEFAULT 0.00,
    "weight" DOUBLE PRECISION,
    "dimensions" VARCHAR(100),
    "shipped_at" TIMESTAMP WITH TIME ZONE,
    "estimated_delivery_at" TIMESTAMP WITH TIME ZONE,
    "delivered_at" TIMESTAMP WITH TIME ZONE,
    "carrier_data" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: suppliers
CREATE TABLE IF NOT EXISTS suppliers (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(200) NOT NULL UNIQUE,
    "code" VARCHAR(50) NOT NULL UNIQUE,
    "contact_name" VARCHAR(200),
    "email" VARCHAR(255),
    "phone" VARCHAR(20),
    "address_line1" VARCHAR(255),
    "address_line2" VARCHAR(255),
    "city" VARCHAR(100),
    "state" VARCHAR(100),
    "postal_code" VARCHAR(20),
    "country" VARCHAR(2),
    "website_url" VARCHAR(500),
    "payment_terms" VARCHAR(100),
    "is_active" BOOLEAN DEFAULT TRUE,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: categories
CREATE TABLE IF NOT EXISTS categories (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "uuid" UUID NOT NULL UNIQUE,
    "parent_id" BIGINT,
    "name" VARCHAR(200) NOT NULL,
    "slug" VARCHAR(250) NOT NULL UNIQUE,
    "description" TEXT,
    "image_url" VARCHAR(500),
    "sort_order" INTEGER DEFAULT 0,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_featured" BOOLEAN DEFAULT FALSE,
    "level" INTEGER DEFAULT 0,
    "path" VARCHAR(500),
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: customers
CREATE TABLE IF NOT EXISTS customers (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "uuid" UUID NOT NULL UNIQUE,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "phone" VARCHAR(20),
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_verified" BOOLEAN DEFAULT FALSE,
    "is_premium" BOOLEAN DEFAULT FALSE,
    "metadata" JSONB,
    "lifetime_value" NUMERIC(12, 2) DEFAULT 0.00,
    "total_orders" BIGINT DEFAULT 0,
    "last_login" TIMESTAMP WITH TIME ZONE,
    "email_verified_at" TIMESTAMP WITH TIME ZONE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: customer_profiles
CREATE TABLE IF NOT EXISTS customer_profiles (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "customer_id" BIGINT NOT NULL UNIQUE,
    "date_of_birth" DATE,
    "gender" VARCHAR(20),
    "avatar_url" VARCHAR(500),
    "bio" TEXT,
    "preferred_language" VARCHAR(10) DEFAULT 'en',
    "timezone" VARCHAR(50) DEFAULT 'UTC',
    "preferences" JSONB,
    "social_links" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: product_variants
CREATE TABLE IF NOT EXISTS product_variants (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "product_id" BIGINT NOT NULL,
    "sku" UUID NOT NULL UNIQUE,
    "name" VARCHAR(200) NOT NULL,
    "option1" VARCHAR(100),
    "option2" VARCHAR(100),
    "option3" VARCHAR(100),
    "price" NUMERIC(12, 2),
    "compare_at_price" NUMERIC(12, 2),
    "stock_quantity" INTEGER DEFAULT 0,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_default" BOOLEAN DEFAULT FALSE,
    "image_url" VARCHAR(500),
    "weight" DOUBLE PRECISION,
    "metadata" JSONB,
    "position" INTEGER DEFAULT 0,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create table: reviews
CREATE TABLE IF NOT EXISTS reviews (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "product_id" BIGINT NOT NULL,
    "customer_id" BIGINT NOT NULL,
    "order_id" BIGINT,
    "rating" INTEGER NOT NULL,
    "title" VARCHAR(200),
    "comment" TEXT,
    "is_verified_purchase" BOOLEAN DEFAULT FALSE,
    "is_approved" BOOLEAN DEFAULT FALSE,
    "helpful_count" INTEGER DEFAULT 0,
    "images" JSONB,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT now(),
    "approved_at" TIMESTAMP WITH TIME ZONE
);

ALTER TABLE addresses ADD CONSTRAINT fk_addresses_customer_id FOREIGN KEY ("customer_id") REFERENCES customers (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE inventory ADD CONSTRAINT fk_inventory_product_id FOREIGN KEY ("product_id") REFERENCES products (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE inventory ADD CONSTRAINT fk_inventory_variant_id FOREIGN KEY ("variant_id") REFERENCES product_variants (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE inventory ADD CONSTRAINT fk_inventory_warehouse_id FOREIGN KEY ("warehouse_id") REFERENCES warehouses (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE order_items ADD CONSTRAINT fk_order_items_order_id FOREIGN KEY ("order_id") REFERENCES orders (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE order_items ADD CONSTRAINT fk_order_items_product_id FOREIGN KEY ("product_id") REFERENCES products (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE order_items ADD CONSTRAINT fk_order_items_variant_id FOREIGN KEY ("variant_id") REFERENCES product_variants (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE payments ADD CONSTRAINT fk_payments_order_id FOREIGN KEY ("order_id") REFERENCES orders (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE products ADD CONSTRAINT fk_products_brand_id FOREIGN KEY ("brand_id") REFERENCES brands (id) ON DELETE SET NULL ON UPDATE NO ACTION;

ALTER TABLE products ADD CONSTRAINT fk_products_supplier_id FOREIGN KEY ("supplier_id") REFERENCES suppliers (id) ON DELETE SET NULL ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_orders_customer_id FOREIGN KEY ("customer_id") REFERENCES customers (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_orders_billing_address_id FOREIGN KEY ("billing_address_id") REFERENCES addresses (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE orders ADD CONSTRAINT fk_orders_shipping_address_id FOREIGN KEY ("shipping_address_id") REFERENCES addresses (id) ON DELETE RESTRICT ON UPDATE NO ACTION;

ALTER TABLE shipping ADD CONSTRAINT fk_shipping_order_id FOREIGN KEY ("order_id") REFERENCES orders (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE categories ADD CONSTRAINT fk_categories_parent_id FOREIGN KEY ("parent_id") REFERENCES categories (id) ON DELETE SET NULL ON UPDATE NO ACTION;

ALTER TABLE customer_profiles ADD CONSTRAINT fk_customer_profiles_customer_id FOREIGN KEY ("customer_id") REFERENCES customers (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE product_variants ADD CONSTRAINT fk_product_variants_product_id FOREIGN KEY ("product_id") REFERENCES products (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_product_id FOREIGN KEY ("product_id") REFERENCES products (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_customer_id FOREIGN KEY ("customer_id") REFERENCES customers (id) ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_order_id FOREIGN KEY ("order_id") REFERENCES orders (id) ON DELETE SET NULL ON UPDATE NO ACTION;

CREATE INDEX IF NOT EXISTS idx_address_customer_id ON addresses ("customer_id");

CREATE INDEX IF NOT EXISTS idx_address_type ON addresses ("type");

CREATE INDEX IF NOT EXISTS idx_address_country ON addresses ("country");

CREATE INDEX IF NOT EXISTS idx_address_postal_code ON addresses ("postal_code");

CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory ("product_id");

CREATE INDEX IF NOT EXISTS idx_inventory_variant_id ON inventory ("variant_id");

CREATE INDEX IF NOT EXISTS idx_inventory_warehouse_id ON inventory ("warehouse_id");

CREATE INDEX IF NOT EXISTS idx_inventory_available_quantity ON inventory ("available_quantity");

CREATE INDEX IF NOT EXISTS idx_order_item_order_id ON order_items ("order_id");

CREATE INDEX IF NOT EXISTS idx_order_item_product_id ON order_items ("product_id");

CREATE INDEX IF NOT EXISTS idx_order_item_variant_id ON order_items ("variant_id");

CREATE INDEX IF NOT EXISTS idx_payment_order_id ON payments ("order_id");

CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_transaction_id ON payments ("transaction_id");

CREATE INDEX IF NOT EXISTS idx_payment_status ON payments ("status");

CREATE INDEX IF NOT EXISTS idx_payment_gateway_transaction_id ON payments ("gateway_transaction_id");

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_sku ON products ("sku");

CREATE UNIQUE INDEX IF NOT EXISTS idx_product_slug ON products ("slug");

CREATE INDEX IF NOT EXISTS idx_product_status ON products ("status");

CREATE INDEX IF NOT EXISTS idx_product_is_active ON products ("is_active");

CREATE INDEX IF NOT EXISTS idx_product_is_featured ON products ("is_featured");

CREATE INDEX IF NOT EXISTS idx_product_brand_id ON products ("brand_id");

CREATE INDEX IF NOT EXISTS idx_product_price ON products ("price");

CREATE UNIQUE INDEX IF NOT EXISTS idx_order_order_number ON orders ("order_number");

CREATE INDEX IF NOT EXISTS idx_order_customer_id ON orders ("customer_id");

CREATE INDEX IF NOT EXISTS idx_order_status ON orders ("status");

CREATE INDEX IF NOT EXISTS idx_order_payment_status ON orders ("payment_status");

CREATE INDEX IF NOT EXISTS idx_order_shipping_status ON orders ("shipping_status");

CREATE INDEX IF NOT EXISTS idx_order_placed_at ON orders ("placed_at");

CREATE UNIQUE INDEX IF NOT EXISTS idx_warehouse_code ON warehouses ("code");

CREATE INDEX IF NOT EXISTS idx_warehouse_is_active ON warehouses ("is_active");

CREATE UNIQUE INDEX IF NOT EXISTS idx_brand_slug ON brands ("slug");

CREATE INDEX IF NOT EXISTS idx_brand_is_active ON brands ("is_active");

CREATE UNIQUE INDEX IF NOT EXISTS idx_shipping_order_id ON shipping ("order_id");

CREATE INDEX IF NOT EXISTS idx_shipping_tracking_number ON shipping ("tracking_number");

CREATE INDEX IF NOT EXISTS idx_shipping_status ON shipping ("status");

CREATE UNIQUE INDEX IF NOT EXISTS idx_supplier_code ON suppliers ("code");

CREATE INDEX IF NOT EXISTS idx_supplier_email ON suppliers ("email");

CREATE INDEX IF NOT EXISTS idx_supplier_is_active ON suppliers ("is_active");

CREATE UNIQUE INDEX IF NOT EXISTS idx_category_slug ON categories ("slug");

CREATE INDEX IF NOT EXISTS idx_category_parent_id ON categories ("parent_id");

CREATE INDEX IF NOT EXISTS idx_category_is_active ON categories ("is_active");

CREATE INDEX IF NOT EXISTS idx_category_level ON categories ("level");

CREATE INDEX IF NOT EXISTS idx_category_path ON categories ("path");

CREATE UNIQUE INDEX IF NOT EXISTS idx_customer_email ON customers ("email");

CREATE UNIQUE INDEX IF NOT EXISTS idx_customer_uuid ON customers ("uuid");

CREATE INDEX IF NOT EXISTS idx_customer_is_active ON customers ("is_active");

CREATE INDEX IF NOT EXISTS idx_customer_created_at ON customers ("created_at");

CREATE UNIQUE INDEX IF NOT EXISTS idx_customer_profile_customer_id ON customer_profiles ("customer_id");

CREATE INDEX IF NOT EXISTS idx_variant_product_id ON product_variants ("product_id");

CREATE UNIQUE INDEX IF NOT EXISTS idx_variant_sku ON product_variants ("sku");

CREATE INDEX IF NOT EXISTS idx_variant_is_active ON product_variants ("is_active");

CREATE INDEX IF NOT EXISTS idx_review_product_id ON reviews ("product_id");

CREATE INDEX IF NOT EXISTS idx_review_customer_id ON reviews ("customer_id");

CREATE INDEX IF NOT EXISTS idx_review_order_id ON reviews ("order_id");

CREATE INDEX IF NOT EXISTS idx_review_rating ON reviews ("rating");

CREATE INDEX IF NOT EXISTS idx_review_is_approved ON reviews ("is_approved");
