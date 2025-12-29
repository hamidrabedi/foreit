# Enterprise Ecommerce Application

This directory contains a complete ecommerce application demonstrating all features of the forge framework, including models, admin interface, REST API, and business logic hooks.

## Models Overview

### Core Models

1. **Customer** - Customer accounts with authentication and profile data
2. **CustomerProfile** - Extended customer profile information (OneToOne with Customer)
3. **Address** - Customer addresses (billing/shipping)
4. **Brand** - Product brands
5. **Supplier** - Product suppliers
6. **Category** - Hierarchical product categories (self-referential)
7. **Product** - Products with variants, inventory, and reviews
8. **ProductVariant** - Product variants (size, color, etc.)
9. **Inventory** - Inventory tracking across warehouses
10. **Warehouse** - Warehouse locations
11. **Order** - Customer orders
12. **OrderItem** - Items within orders
13. **Payment** - Payment records
14. **Shipping** - Shipping information (OneToOne with Order)
15. **Review** - Product reviews

## PostgreSQL Data Types Used

### Numeric Types
- `BIGINT` - IDs, counts (Int64)
- `INTEGER` - Smaller integers (Int32)
- `NUMERIC(precision, scale)` - Decimal fields (Decimal with MaxDigits/DecimalPlaces)
- `DOUBLE PRECISION` - Float64 fields
- `REAL` - Float32 fields

### String Types
- `VARCHAR(n)` - String fields with MaxLength
- `TEXT` - Text fields without length limit
- `CHAR(36)` - UUID fields (or UUID type in PostgreSQL)

### Boolean
- `BOOLEAN` - Bool fields

### Date/Time Types
- `DATE` - Date fields
- `TIME` - Time fields
- `TIMESTAMP WITH TIME ZONE` - DateTime fields
- `TIMESTAMP` - Time fields (fallback)

### Special Types
- `UUID` - UUID fields
- `JSONB` - JSON fields
- `BYTEA` - Bytes fields

## Relations Demonstrated

### ForeignKey Relations
- Customer → Address (billing_address_id, shipping_address_id)
- Customer → CustomerProfile (OneToOne)
- Address → Customer
- Product → Brand
- Product → Supplier
- Product → Category (ManyToMany)
- ProductVariant → Product
- Inventory → Product, ProductVariant, Warehouse
- Order → Customer, Address (billing, shipping)
- OrderItem → Order, Product, ProductVariant
- Payment → Order
- Shipping → Order (OneToOne)
- Review → Product, Customer, Order
- Category → Category (self-referential, parent_id)

### OneToOne Relations
- Customer ↔ CustomerProfile
- Order ↔ Shipping

### ManyToMany Relations
- Product ↔ Category (through junction table)

### OneToMany Relations
- Customer → Addresses
- Customer → Orders
- Customer → Reviews
- Product → Variants
- Product → Inventory Items
- Product → Reviews
- Order → OrderItems
- Order → Payments
- Brand → Products
- Supplier → Products
- Category → Children (self-referential)
- Warehouse → Inventory Items

## Field Options Used

- **Primary Keys**: Auto-incrementing Int64 IDs
- **UUIDs**: Unique identifiers
- **Required/Optional**: Field nullability
- **Unique**: Unique constraints
- **MaxLength**: String length limits
- **MaxDigits/DecimalPlaces**: Decimal precision
- **Default Values**: Various types
- **Choices**: Enum-like string fields
- **AutoNow/AutoNowAdd**: Automatic timestamp management
- **DBIndex**: Database indexes
- **JSON**: Complex nested data
- **WriteOnly**: Fields not serialized (password_hash)

## Indexes

Each model includes appropriate indexes for:
- Foreign keys
- Frequently queried fields
- Unique constraints
- Composite indexes (UniqueTogether)

## Cascade Behaviors

- **CASCADE**: Delete related records when parent is deleted
- **SET_NULL**: Set foreign key to NULL when parent is deleted
- **PROTECT**: Prevent deletion if related records exist

## Application Structure

```
ecommerce/
├── cmd/
│   └── server/
│       └── main.go          # Main application entry point
├── api/
│   ├── serializers.go       # REST API serializers
│   └── viewsets.go          # REST API viewsets
├── models/                  # Model definitions
│   ├── *.go                 # Schema definitions
│   └── *.gen.go             # Generated code (after forge generate)
├── migrations/              # Database migrations
├── config/
│   └── config.yaml          # Configuration
└── README.md

```

## Quick Start

1. **Generate Code**: `forge generate`
2. **Run Migrations**: `forge migrate`
3. **Start Server**: `go run cmd/server/main.go`
4. **Access Admin**: http://localhost:8000/admin/
5. **Access API**: http://localhost:8000/api/v1/

See [SETUP.md](SETUP.md) for detailed setup instructions.

## Usage

These models demonstrate:
1. All supported PostgreSQL data types
2. Complex nested relations (3+ levels deep)
3. Self-referential relations (Category)
4. ManyToMany relations
5. OneToOne relations
6. Composite unique constraints
7. Complex indexes
8. JSON fields for flexible data storage
9. Decimal precision handling
10. UUID primary/foreign keys

## Testing the Schema

This example demonstrates all SQL features supported by the ForgeGo ORM. The models are designed to test:

### Running Tests

The test suite automatically generates code using the `forge` CLI tool before running verification tests. Make sure the `forge` binary is built and available.

To verify that all SQL features are properly supported:

```bash
# First, build the forge CLI tool (from newforge directory)
cd ../../..
go build -o forge.exe ./cmd/forge  # or ./forge on Unix

# Then run tests (tests will auto-generate code first)
cd examples/ecommerce
go test -v ./models

# Or run specific test suites
go test -v ./models -run TestSQLFeatureSupport
go test -v ./models -run TestPostgreSQLTypes
go test -v ./models -run TestComplexRelations
go test -v ./models -run TestFieldOptions
go test -v ./models -run TestIndexesAndConstraints
```

**Note**: Tests will automatically run `forge generate` before verification, so you don't need to generate code manually.

### Creating Migrations

Since we now use golang-migrate for migrations, you create migration files manually:

```bash
# Create a new migration file
forge makemigrations create_ecommerce_tables --path ./migrations

# Then manually write SQL in the generated .up.sql and .down.sql files
# Apply migrations
forge migrate --path ./migrations
```

### What's Included

- **Model Builders**: All 15 model definitions in `models/*.go` (original source files)
- **Generated Code**: Auto-generated by tests using `forge generate` command (not committed)
- **Tests**: Comprehensive test suite in `models/sql_verification_test.go` that:
  - Automatically generates code using the `forge` CLI tool
  - Verifies all SQL features are properly supported:
  - All PostgreSQL data types are supported
  - All field options work correctly
  - All relation types are properly defined
  - Indexes and constraints are correctly specified
  - Cascade behaviors are supported

### Model Files

- `address.go` - Customer addresses
- `brand.go` - Product brands
- `category.go` - Hierarchical categories (self-referential)
- `customer.go` - Customer accounts
- `customer_profile.go` - Extended customer profiles (OneToOne)
- `inventory.go` - Inventory tracking
- `order.go` - Customer orders
- `order_item.go` - Order line items
- `payment.go` - Payment records
- `product.go` - Products with variants
- `product_variant.go` - Product variants
- `review.go` - Product reviews
- `shipping.go` - Shipping information (OneToOne with Order)
- `supplier.go` - Product suppliers
- `warehouse.go` - Warehouse locations

