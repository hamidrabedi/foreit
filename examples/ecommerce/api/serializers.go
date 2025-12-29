package api

import (
	"github.com/forgego/forge/pkg/api"
	"ecommerce/models"
)

// CustomerSerializer serializes Customer model
type CustomerSerializer struct {
	*api.BaseSerializer
}

func NewCustomerSerializer() api.Serializer {
	return &CustomerSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *CustomerSerializer) Fields() []string {
	return []string{
		"id", "uuid", "email", "phone", "first_name", "last_name",
		"is_active", "is_verified", "is_premium", "lifetime_value",
		"total_orders", "last_login", "email_verified_at",
		"created_at", "updated_at",
	}
}

func (s *CustomerSerializer) ReadOnlyFields() []string {
	return []string{"id", "uuid", "created_at", "updated_at", "last_login", "email_verified_at"}
}

// ProductSerializer serializes Product model
type ProductSerializer struct {
	*api.BaseSerializer
}

func NewProductSerializer() api.Serializer {
	return &ProductSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *ProductSerializer) Fields() []string {
	return []string{
		"id", "sku", "name", "slug", "description", "short_description",
		"price", "compare_at_price", "cost_price", "currency",
		"stock_quantity", "low_stock_threshold", "track_inventory",
		"is_active", "is_featured", "is_digital", "requires_shipping",
		"weight", "length", "width", "height", "status",
		"brand_id", "supplier_id", "published_at",
		"created_at", "updated_at",
	}
}

func (s *ProductSerializer) ReadOnlyFields() []string {
	return []string{"id", "sku", "created_at", "updated_at"}
}

// OrderSerializer serializes Order model
type OrderSerializer struct {
	*api.BaseSerializer
}

func NewOrderSerializer() api.Serializer {
	return &OrderSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *OrderSerializer) Fields() []string {
	return []string{
		"id", "order_number", "customer_id", "status",
		"subtotal", "tax_amount", "shipping_amount", "discount_amount",
		"total_amount", "currency", "billing_address_id", "shipping_address_id",
		"payment_status", "shipping_status", "notes", "customer_notes",
		"placed_at", "confirmed_at", "shipped_at", "delivered_at",
		"cancelled_at", "created_at", "updated_at",
	}
}

func (s *OrderSerializer) ReadOnlyFields() []string {
	return []string{"id", "order_number", "placed_at", "created_at", "updated_at"}
}

// OrderItemSerializer serializes OrderItem model
type OrderItemSerializer struct {
	*api.BaseSerializer
}

func NewOrderItemSerializer() api.Serializer {
	return &OrderItemSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *OrderItemSerializer) Fields() []string {
	return []string{
		"id", "order_id", "product_id", "product_variant_id",
		"quantity", "unit_price", "total_price", "discount_amount",
		"created_at", "updated_at",
	}
}

func (s *OrderItemSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// AddressSerializer serializes Address model
type AddressSerializer struct {
	*api.BaseSerializer
}

func NewAddressSerializer() api.Serializer {
	return &AddressSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *AddressSerializer) Fields() []string {
	return []string{
		"id", "customer_id", "type", "street_address", "city",
		"state", "postal_code", "country", "is_default",
		"created_at", "updated_at",
	}
}

func (s *AddressSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// BrandSerializer serializes Brand model
type BrandSerializer struct {
	*api.BaseSerializer
}

func NewBrandSerializer() api.Serializer {
	return &BrandSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *BrandSerializer) Fields() []string {
	return []string{
		"id", "name", "slug", "description", "website", "logo_url",
		"is_active", "created_at", "updated_at",
	}
}

func (s *BrandSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// SupplierSerializer serializes Supplier model
type SupplierSerializer struct {
	*api.BaseSerializer
}

func NewSupplierSerializer() api.Serializer {
	return &SupplierSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *SupplierSerializer) Fields() []string {
	return []string{
		"id", "name", "contact_name", "email", "phone", "address",
		"city", "state", "postal_code", "country", "is_active",
		"created_at", "updated_at",
	}
}

func (s *SupplierSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// CategorySerializer serializes Category model
type CategorySerializer struct {
	*api.BaseSerializer
}

func NewCategorySerializer() api.Serializer {
	return &CategorySerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *CategorySerializer) Fields() []string {
	return []string{
		"id", "name", "slug", "description", "parent_id",
		"is_active", "created_at", "updated_at",
	}
}

func (s *CategorySerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// ProductVariantSerializer serializes ProductVariant model
type ProductVariantSerializer struct {
	*api.BaseSerializer
}

func NewProductVariantSerializer() api.Serializer {
	return &ProductVariantSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *ProductVariantSerializer) Fields() []string {
	return []string{
		"id", "product_id", "sku", "name", "price", "cost_price",
		"stock_quantity", "is_active", "attributes",
		"created_at", "updated_at",
	}
}

func (s *ProductVariantSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// InventorySerializer serializes Inventory model
type InventorySerializer struct {
	*api.BaseSerializer
}

func NewInventorySerializer() api.Serializer {
	return &InventorySerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *InventorySerializer) Fields() []string {
	return []string{
		"id", "product_id", "product_variant_id", "warehouse_id",
		"quantity", "reserved_quantity", "available_quantity",
		"created_at", "updated_at",
	}
}

func (s *InventorySerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// WarehouseSerializer serializes Warehouse model
type WarehouseSerializer struct {
	*api.BaseSerializer
}

func NewWarehouseSerializer() api.Serializer {
	return &WarehouseSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *WarehouseSerializer) Fields() []string {
	return []string{
		"id", "name", "code", "address", "city", "state",
		"postal_code", "country", "is_active",
		"created_at", "updated_at",
	}
}

func (s *WarehouseSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// PaymentSerializer serializes Payment model
type PaymentSerializer struct {
	*api.BaseSerializer
}

func NewPaymentSerializer() api.Serializer {
	return &PaymentSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *PaymentSerializer) Fields() []string {
	return []string{
		"id", "order_id", "payment_method", "amount", "currency",
		"status", "transaction_id", "processed_at",
		"created_at", "updated_at",
	}
}

func (s *PaymentSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// ShippingSerializer serializes Shipping model
type ShippingSerializer struct {
	*api.BaseSerializer
}

func NewShippingSerializer() api.Serializer {
	return &ShippingSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *ShippingSerializer) Fields() []string {
	return []string{
		"id", "order_id", "carrier", "tracking_number", "status",
		"shipped_at", "delivered_at", "created_at", "updated_at",
	}
}

func (s *ShippingSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}

// ReviewSerializer serializes Review model
type ReviewSerializer struct {
	*api.BaseSerializer
}

func NewReviewSerializer() api.Serializer {
	return &ReviewSerializer{
		BaseSerializer: api.NewBaseSerializer(),
	}
}

func (s *ReviewSerializer) Fields() []string {
	return []string{
		"id", "product_id", "customer_id", "order_id",
		"rating", "title", "comment", "is_verified_purchase",
		"is_approved", "created_at", "updated_at",
	}
}

func (s *ReviewSerializer) ReadOnlyFields() []string {
	return []string{"id", "created_at", "updated_at"}
}
