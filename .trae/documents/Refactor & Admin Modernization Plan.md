# Ultra Modernization & Admin Overhaul Plan

## 1. Code Cleanup & Refactoring ("The Purge")

* **Legacy Elimination**:

  * Migrate remaining `ecommerce` models (Inventory, Marketing, Orders, Customers) to Functional Schema API.

  * Delete legacy Schema Builder files (`typed_builders.go`, `unified_field_builder.go`).

  * Grep and remove all `@deprecated` code and compatibility layers.

* **Quality & Naming**:

  * Audit `forge/server` and `forge/api` for naming conventions.

  * Refactor duplicate logic in `forge/orm` and `forge/server` into shared helpers.

  * Remove redundant/low-value documentation.

## 2. Admin System (Backend)

* **Architecture**:

  * Implement `Plugin` interface in `forge/admin/core`.

  * Create `Registry` for Widgets, Custom Pages, and Menu Items.

* **API**:

  * Implement `GET /admin/api/config` to serve dynamic layout/menu configuration.

  * Ensure Admin API automatically uses `forge/validate` for all inputs.

* **Type Safety**:

  * Refactor `admin.Register` to enforce strict typing against Schema definitions.

## 3. Admin UI (Frontend)

* **Modernization**:

  * Refactor `forge/admin/ui/web` to use **TanStack Query** (data) and **Zustand** (state).

* **Extensibility**:

  * Implement `ComponentRegistry` (Frontend) to map API keys to React components.

  * Create **Dynamic Dashboard** supporting custom widgets.

  * Create **Dynamic Sidebar** driven by Backend API.

* **Features**:

  * Add support for Custom Pages and Overriding built-in views.

## 4. Server & Data Layer

* **Filtering**:

  * Implement type-safe `Filter` system in ORM (`forge/orm/filters`).

  * Auto-generate Filter UI in Admin based on Schema types.

* **Server**:

  * Modularize `forge/server` middleware and routing.

## 5. Verification

* **Build**: Ensure `examples/ecommerce` compiles with all changes.

* **Test**: Verify Admin UI loads and renders dynamic configuration.

