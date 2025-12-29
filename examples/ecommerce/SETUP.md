# Ecommerce Example - Setup Guide

This is a comprehensive ecommerce application demonstrating all features of the forge framework.

## Features Implemented

### ✅ Framework Features Used

1. **Schema Definition System**
   - All 15 models with complete field definitions
   - All PostgreSQL data types (Int64, String, Text, Decimal, Bool, DateTime, UUID, JSON)
   - Field options (Required, Optional, Unique, Default, Choices, AutoNow, AutoNowAdd)
   - Indexes and constraints
   - UniqueTogether constraints

2. **Relations**
   - ForeignKey relations
   - OneToOne relations (Customer ↔ CustomerProfile, Order ↔ Shipping)
   - OneToMany relations (Customer → Orders, Product → Variants)
   - ManyToMany relations (Product ↔ Category)
   - Self-referential relations (Category → Category)

3. **Model Hooks**
   - Order: Total calculation hooks
   - OrderItem: Price calculation and inventory hooks
   - Customer: Password hashing hooks (structure ready)

4. **Admin Interface**
   - All models registered with admin
   - Custom list displays
   - Search fields configured
   - List filters configured
   - Auto-generated forms

5. **REST API**
   - Serializers for all models
   - ViewSets for all models
   - Full CRUD operations
   - Pagination, filtering, ordering support

6. **Server Setup**
   - Database connection
   - Admin routes
   - REST API routes
   - Health check endpoint

## Setup Instructions

### 1. Generate Code

First, generate the model code from schema definitions:

```bash
# From the project root
cd examples/ecommerce
forge generate
```

This will create:
- `models/*.gen.go` - Generated model structs
- `models/*_manager.gen.go` - Generated managers
- `models/*_fields.gen.go` - Generated field expressions
- `models/*_queryset.gen.go` - Generated querysets

### 2. Update main.go

After code generation, update `cmd/server/main.go` to uncomment the manager setup:

```go
func setManagersDB(database *forge.DB) {
	// Uncomment these after code generation:
	models.Customer.Objects.SetDB(database)
	models.CustomerProfile.Objects.SetDB(database)
	models.Address.Objects.SetDB(database)
	models.Brand.Objects.SetDB(database)
	models.Supplier.Objects.SetDB(database)
	models.Category.Objects.SetDB(database)
	models.Product.Objects.SetDB(database)
	models.ProductVariant.Objects.SetDB(database)
	models.Inventory.Objects.SetDB(database)
	models.Warehouse.Objects.SetDB(database)
	models.Order.Objects.SetDB(database)
	models.OrderItem.Objects.SetDB(database)
	models.Payment.Objects.SetDB(database)
	models.Shipping.Objects.SetDB(database)
	models.Review.Objects.SetDB(database)
}
```

### 3. Update Admin Registration

After code generation, update admin registration to include managers:

```go
func registerAdminModels() {
	// Use RegisterModelWithManager after code generation:
	admin.RegisterModelWithManager(&models.Customer{}, models.Customer.Objects)
	// ... etc for all models
}
```

### 4. Update API ViewSets

After code generation, update viewsets to use actual querysets:

```go
func registerCustomerAPI(apiRouter *api.Router) {
	viewset := api.NewBaseViewSet(
		NewCustomerSerializer,
		models.Customer.Objects.Filter(), // Use generated queryset
		&models.Customer{},
	)
	apiRouter.Register("customers", viewset)
}
```

### 5. Run Migrations

```bash
forge migrate
```

### 6. Start Server

```bash
go run cmd/server/main.go
```

Or use the forge CLI:

```bash
forge runserver
```

## Access Points

- **Admin Interface**: http://localhost:8000/admin/
- **REST API**: http://localhost:8000/api/v1/
- **Health Check**: http://localhost:8000/

## API Endpoints

All models have full CRUD endpoints:

- `GET /api/v1/customers/` - List customers
- `POST /api/v1/customers/` - Create customer
- `GET /api/v1/customers/{id}/` - Get customer
- `PATCH /api/v1/customers/{id}/` - Update customer
- `DELETE /api/v1/customers/{id}/` - Delete customer

Same pattern for:
- `/api/v1/products/`
- `/api/v1/orders/`
- `/api/v1/order-items/`
- `/api/v1/categories/`
- `/api/v1/brands/`
- `/api/v1/suppliers/`
- `/api/v1/inventory/`
- `/api/v1/warehouses/`
- `/api/v1/payments/`
- `/api/v1/shipping/`
- `/api/v1/reviews/`
- `/api/v1/addresses/`

## Model Hooks

### Order Hooks
- `BeforeSave`: Calculates total_amount from subtotal + tax + shipping - discount
- `BeforeCreate`: Sets placed_at timestamp
- `AfterCreate`: Placeholder for notifications/logging

### OrderItem Hooks
- `BeforeSave`: Calculates total_price from unit_price * quantity
- `BeforeCreate`: Validates inventory availability (structure ready)
- `AfterCreate`: Updates order totals (structure ready)
- `BeforeDelete`: Releases inventory (structure ready)

### Customer Hooks
- `BeforeCreate`: Password hashing (structure ready)
- `BeforeUpdate`: Password hashing on update (structure ready)

## Admin Features

### List Views
- Pagination
- Search (on configured fields)
- Filtering (on configured fields)
- Sorting (click column headers)

### Detail Views
- View all fields
- Related objects (coming soon)

### Forms
- Auto-generated create/edit forms
- Field validation
- Read-only fields

## Next Steps

1. **Complete Hooks**: Implement inventory updates, order total recalculation
2. **Add Authentication**: Implement customer authentication
3. **Add Permissions**: Add role-based access control
4. **Add Caching**: Cache frequently accessed data
5. **Add Background Tasks**: Process orders, send emails
6. **Add Webhooks**: Notify external systems of events

## Testing

Run the test suite:

```bash
go test ./models -v
```

This will:
- Generate code automatically
- Verify all SQL features
- Test all data types
- Test all relations
- Test all field options
