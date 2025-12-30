# forge Framework - Code Structure Reference

This document provides a comprehensive overview of the forge framework codebase structure, explaining what each package, file, and folder does.

## Table of Contents

- [Root Directory Structure](#root-directory-structure)
- [Package Structure (`pkg/`)](#package-structure-pkg)
  - [Core Framework Packages](#core-framework-packages)
  - [Feature Packages](#feature-packages)
  - [Infrastructure Packages](#infrastructure-packages)
- [Command Structure (`cmd/`)](#command-structure-cmd)
- [Examples (`examples/`)](#examples-examples)
- [Tests (`tests/`)](#tests-tests)
- [Documentation (`docs/`)](#documentation-docs)

---

## Root Directory Structure

```
foreit/
├── cmd/                    # Application entry points
├── pkg/                    # Framework packages (main codebase)
├── docs/                   # Documentation
├── examples/              # Example applications
├── tests/                 # Test packages
├── helper-projects/        # Reference implementations (Django, etc.)
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── README.md               # Project overview
├── Makefile                # Build automation
└── config.go, db.go, etc.  # Root-level utilities (legacy)
```

---

## Package Structure (`pkg/`)

The `pkg/` directory contains all framework packages. Each package is self-contained with a specific responsibility.

### Core Framework Packages

#### `pkg/schema/` - Schema Definition System

**Purpose:** Declarative model definitions with Django-like field options.

**Files:**
- `schema.go` - Core Schema interface and BaseSchema
- `field.go` - Field type definitions and field structure
- `field_builder.go` - Field builder implementations
- `typed_builders.go` - Type-specific field builders (Int64, String, Bool, etc.)
- `unified_field_builder.go` - Unified field builder API
- `field_traits.go` - Field trait mixins (Required, Unique, etc.)
- `field_config.go` - Field configuration options
- `relation.go` - Relationship definitions (ForeignKey, OneToOne, ManyToMany)
- `meta.go` - Model metadata (Meta options)
- `hooks.go` - Lifecycle hooks (BeforeSave, AfterCreate, etc.)
- `registry.go` - Schema registry for model discovery
- `constraint_builder.go` - Database constraint builders
- `index_builder.go` - Database index builders
- `helpers.go` - Schema utility functions

**Key Types:**
- `Schema` - Interface that models implement
- `Field` - Field definition structure
- `Relation` - Relationship definition
- `Meta` - Model metadata
- `ModelHooks` - Lifecycle hooks

---

#### `pkg/generator/` - Code Generation System

**Purpose:** AST-based code generation from schema definitions.

**Files:**
- `generator.go` - Main generator orchestrator
- `ast_parser.go` - Go AST parser for extracting schema definitions
- `writer.go` - Code file writer
- `templates.go` - Template management
- `templates/` - Code generation templates
  - `manager.tmpl` - Manager generation template
  - `queryset.tmpl` - QuerySet generation template
  - `fields.tmpl` - FieldExpr generation template
  - `model.tmpl` - Model struct generation template

**Key Types:**
- `Generator` - Main generator struct
- `ASTParser` - AST parsing logic
- `Writer` - Code file writing
- `ModelDefinition` - Parsed model definition

**Generated Files:**
- `models/*.gen.go` - Generated model structs
- `models/*_fields.gen.go` - Generated FieldExpr definitions
- `models/*_manager.gen.go` - Generated Manager with CRUD
- `models/*_queryset.gen.go` - Generated QuerySet wrappers

---

#### `pkg/query/` - Type-Safe ORM

**Purpose:** Type-safe query building and execution with Django-like API.

**Files:**
- `queryset.go` - Base QuerySet interface and implementation
- `queryset.go` - QuerySet with additional features
- `manager.go` - Manager interface and base implementation
- `manager.go` - Manager with CRUD operations
- `field_expr.go` - Type-safe field expressions (`FieldExpr[T]`)
- `query_expr.go` - Query expression building (`QueryExpr`)
- `expression.go` - Expression evaluation and SQL generation
- `sql_builder.go` - SQL query builder with parameter binding
- `aggregates.go` - Aggregate functions (Count, Sum, Avg, etc.)
- `annotations.go` - Annotation support (computed fields)
- `update_builder.go` - Update query builder
- `related_field.go` - Related field handling
- `schema.go` - QuerySet schema integration
- `types.go` - Type definitions
- `field_helpers.go` - Field expression helpers
- `manager_helpers.go` - Manager helper functions
- `update_helpers.go` - Update operation helpers
- `update_builder_helpers.go` - Update builder helpers

**Key Types:**
- `QuerySet[T]` - Type-safe query set interface
- `BaseQuerySet[T]` - Base QuerySet implementation
- `Manager[T]` - Manager interface
- `FieldExpr[T, F]` - Type-safe field expression
- `QueryExpr` - Query expression
- `Aggregate` - Aggregate function
- `AnnotationExpr` - Annotation expression

---

#### `pkg/db/` - Database Layer

**Purpose:** Database connection, transactions, and migration integration.

**Files:**
- `db.go` - Database connection wrapper
- `transaction.go` - Transaction management
- `migrations.go` - Migration system integration

**Key Types:**
- `DB` - Database connection wrapper
- `Config` - Database configuration
- `Tx` - Transaction wrapper

---

#### `pkg/migrations/` - Migration System

**Purpose:** Complete migration system with state management, detection, and SQL generation.

**Files:**
- `generator.go` - Migration generator
- `detector.go` - Schema change detector
- `state.go` - Migration state management
- `plan.go` - Migration planning
- `sqlgen.go` - SQL generation entry point
- `linter.go` - Migration linter
- `changes.go` - Change tracking

**Subdirectories:**
- `core/` - Core migration types and interfaces
- `detection/` - Schema change detection
- `generation/` - Migration generation
- `execution/` - Migration execution
- `state/` - State management and parsing
- `sql/` - SQL generation (PostgreSQL, SQLite)
- `operations/` - Django-style operations
- `validation/` - Migration validation
- `recovery/` - Migration recovery
- `squashing/` - Migration squashing
- `dependencies/` - Dependency resolution
- `linting/` - Migration linting
- `drift/` - Schema drift detection
- `logging/` - Migration logging
- `utils/` - Utility functions

---

#### `pkg/migrate/` - Migration Implementation

**Purpose:** Lower-level migration implementation and SQL parsing.

**Files:**
- `migrate.go` - Migration interface
- `migrate_impl.go` - Migration implementation

**Subdirectories:**
- `generate/` - Migration generation
- `apply/` - Migration application
- `verify/` - Migration verification
- `schema/` - Schema loading and conversion
- `sqlgen/` - SQL generation
- `sqlparse/` - SQL parsing
- `types/` - Migration types
- `errors/` - Migration errors
- `internal/` - Internal utilities

---

### Feature Packages

#### `pkg/admin/` - Admin System

**Purpose:** Type-safe admin interface with full CRUD operations.

**Files:**
- `admin.go` - Core Admin[T] and Config[T] types
- `registry.go` - Admin model registry
- `list_view.go` - List view implementation
- `detail_view.go` - Detail view implementation
- `form_view.go` - Create/update form views
- `form.go` - Form handling and validation
- `fields.go` - Type-safe field expressions for admin
- `filters.go` - Filter system (Boolean, Choice)
- `actions.go` - Bulk actions
- `widgets.go` - Form widgets
- `export.go` - CSV/JSON export
- `inlines.go` - Inline editing support
- `fieldsets.go` - Form field grouping
- `ordering.go` - Ordering system
- `permissions.go` - Permission checking
- `validation.go` - Form validation
- `errors.go` - Admin-specific errors
- `example.go` - Usage examples
- `complete_example.go` - Complete production example
- `registry_export.go` - Registry export utilities

**Subdirectories:**
- `http/` - HTTP handlers and routing
  - `handlers.go` - HTTP request handlers
  - `router.go` - Admin route registration
  - `type_registry.go` - Type registry for HTTP layer
  - `init.go` - Initialization
  - `handlers_test.go` - Handler tests
- `testing/` - Testing utilities
  - `helpers.go` - Test helpers

**Key Types:**
- `Admin[T]` - Type-safe admin instance
- `Config[T]` - Admin configuration
- `FieldExpr[T, F]` - Type-safe field expression
- `Filter[T]` - Filter interface
- `Action[T]` - Bulk action
- `ListView[T]` - List view
- `DetailView[T]` - Detail view
- `FormView[T]` - Form view

---

#### `pkg/api/` - REST API Framework

**Purpose:** DRF-like REST API framework with serializers, viewsets, authentication, and permissions.

**Files:**
- `api.go` - API framework settings and initialization
- `viewset.go` - Base ViewSet interface and implementation
- `viewset_enhanced.go` - Enhanced ViewSet with additional features
- `viewset_enhanced_integrated.go` - Integrated enhanced ViewSet
- `serializer.go` - Base serializer interface
- `router_enhanced.go` - Enhanced router with action support
- `pagination.go` - Pagination (PageNumber, LimitOffset)
- `filters.go` - Filtering system
- `content_negotiation.go` - Content negotiation
- `integration.go` - API integration helpers
- `helpers.go` - API helper functions
- `middleware_integration.go` - Middleware integration

**Subdirectories:**
- `core/` - Core framework components
  - `request.go` - Request wrapper
  - `response.go` - Response wrapper
  - `context.go` - Request context utilities
  - `middleware.go` - Core middleware interface
- `serializers/` - Serialization system
  - `serializer.go` - Base serializer
  - `enhanced_serializer.go` - Enhanced serializer
  - `fields/` - Serializer field types
    - `base.go` - Base field
    - `string.go` - StringField
    - `integer.go` - IntegerField
    - `float.go` - FloatField
    - `boolean.go` - BooleanField
    - `datetime.go` - DateTimeField
    - `email.go` - EmailField
    - `url.go` - URLField
    - `uuid.go` - UUIDField
    - `readonly.go` - ReadOnlyField
    - `writeonly.go` - WriteOnlyField
    - `hidden.go` - HiddenField
    - `method.go` - SerializerMethodField
- `authentication/` - Authentication system
  - `auth.go` - Base authentication interface
  - `token.go` - TokenAuthentication
  - `jwt.go` - JWTAuthentication
  - `session.go` - SessionAuthentication
  - `basic.go` - BasicAuthentication
  - `apikey.go` - APIKeyAuthentication
  - `result.go` - Authentication result
- `permissions/` - Permission system
  - `permission.go` - Base permission interface
  - `allow_any.go` - AllowAny
  - `is_authenticated.go` - IsAuthenticated
  - `is_authenticated_or_read_only.go` - IsAuthenticatedOrReadOnly
  - `is_admin.go` - IsAdminUser
  - `is_owner.go` - IsOwnerOrReadOnly
  - `helpers.go` - Permission helpers
- `throttling/` - Throttling system
  - `throttle.go` - Base throttle interface
  - `anon_rate.go` - AnonRateThrottle
  - `user_rate.go` - UserRateThrottle
  - `scoped_rate.go` - ScopedRateThrottle
  - `cache.go` - Throttle cache backend
  - `rate_parser.go` - Rate string parser
- `renderers/` - Response renderers
  - `renderer.go` - Base renderer interface
  - `json.go` - JSONRenderer
  - `xml.go` - XMLRenderer
  - `yaml.go` - YAMLRenderer
  - `html.go` - HTMLRenderer
  - `csv.go` - CSVRenderer
- `parsers/` - Request parsers
  - `parser.go` - Base parser interface
  - `json.go` - JSONParser
  - `xml.go` - XMLParser
  - `form.go` - FormParser
  - `multipart.go` - MultiPartParser
  - `yaml.go` - YAMLParser
- `filters/` - Filtering system
  - `backend.go` - Filter backend interface
  - `search.go` - SearchFilter
  - `ordering.go` - OrderingFilter
- `exceptions/` - Exception handling
  - `exception.go` - Base exception
  - `validation.go` - ValidationError
  - `auth.go` - Authentication exceptions
  - `permission.go` - Permission exceptions
  - `not_found.go` - NotFound
  - `throttled.go` - Throttled
  - `handler.go` - Exception handler
- `versioning/` - API versioning
  - `version.go` - Version interface
  - `url_path.go` - URLPathVersioning
  - `query_param.go` - QueryParameterVersioning
  - `header.go` - HeaderVersioning
- `caching/` - Caching system
  - `backend.go` - Cache backend interface
  - `memory.go` - MemoryCache
- `docs/` - Documentation generation
  - `openapi.go` - OpenAPI generator
- `testing/` - Testing utilities
  - `client.go` - APIClient

**Key Types:**
- `ViewSet` - ViewSet interface
- `BaseViewSet` - Base ViewSet implementation
- `Serializer` - Serializer interface
- `Authentication` - Authentication interface
- `Permission` - Permission interface
- `Throttle` - Throttle interface
- `Renderer` - Renderer interface
- `Parser` - Parser interface

---

#### `pkg/users/` - User System

**Purpose:** Complete user management and authentication system.

**Files:**
- `factory.go` - User system factory
- `router.go` - User system router setup

**Subdirectories:**
- `models/` - User models
  - `user.go` - User model
  - `session.go` - UserSession model
  - `permission.go` - Permission model
  - `token.go` - Token models (EmailVerification, PasswordReset)
- `repository/` - Data access layer
  - `interface.go` - Repository interfaces
  - `user.go` - User repository implementation
  - `session.go` - Session repository
  - `permission.go` - Permission repository
  - `token.go` - Token repository
  - `user_test.go` - User repository tests
- `service/` - Business logic layer
  - `interface.go` - Service interfaces
  - `user.go` - User service
  - `auth.go` - Authentication service
  - `password.go` - Password service
  - `permission.go` - Permission service
  - `user_test.go` - User service tests
  - `auth_test.go` - Auth service tests
  - `password_test.go` - Password service tests
- `backends/` - Authentication backends
  - `interface.go` - Backend interface
  - `password.go` - Password backend
  - `token.go` - Token backend
  - `registry.go` - Backend registry
  - `password_test.go` - Password backend tests
- `serializers/` - API serializers
  - `user.go` - User serializers
  - `auth.go` - Auth serializers
- `handlers/` - HTTP handlers
  - `user.go` - User viewset
  - `auth.go` - Auth endpoints
- `middleware/` - Authentication middleware
  - `auth.go` - Auth middleware
- `config/` - Configuration
  - `settings.go` - User system settings
- `migrations/` - Database migrations
  - `0001_create_users.up.sql` - User tables migration
  - `0001_create_users.down.sql` - Rollback migration

**Key Types:**
- `User` - User model
- `UserSession` - Session model
- `Permission` - Permission model
- `Group` - Group model
- `UserRepository` - User repository interface
- `UserService` - User service interface
- `AuthService` - Authentication service interface
- `AuthenticationBackend` - Backend interface

---

### Infrastructure Packages

#### `pkg/http/` - HTTP & Routing

**Purpose:** HTTP server, routing, and middleware.

**Files:**
- `router.go` - Chi router wrapper
- `middleware.go` - Middleware utilities
- `context.go` - Request context utilities
- `server.go` - Server wrapper
- `response.go` - Response utilities
- `request.go` - Request utilities

**Key Types:**
- `Router` - Framework router wrapper
- `Server` - Server wrapper

---

#### `pkg/config/` - Configuration

**Purpose:** Application configuration management.

**Files:**
- `config.go` - Viper wrapper and configuration loader
- `settings.go` - Framework settings structure

**Key Types:**
- `Config` - Configuration wrapper
- `Settings` - Framework settings

---

#### `pkg/logging/` - Logging

**Purpose:** Structured logging with zap.

**Files:**
- `logger.go` - Zap logger wrapper
- `middleware.go` - Request logging middleware

**Key Types:**
- `Logger` - Logger wrapper

---

#### `pkg/security/` - Security

**Purpose:** Security features (CSRF, XSS, SQL injection prevention).

**Files:**
- `csrf.go` - CSRF protection
- `sessions.go` - Session management
- `xss.go` - XSS protection
- `sql_injection.go` - SQL injection prevention

---

#### `pkg/validation/` - Validation

**Purpose:** Data validation with go-playground/validator.

**Files:**
- `validator.go` - Validator wrapper
- `tags.go` - Validation tag helpers
- `integration.go` - Schema integration
- `field_validator.go` - Field-level validation
- `helpers.go` - Validation helpers

**Key Types:**
- `Validator` - Validator wrapper

---

#### `pkg/auth/` - Basic Authentication

**Purpose:** Basic authentication utilities.

**Files:**
- `password.go` - Password hashing (bcrypt)

---

#### `pkg/errors/` - Error Handling

**Purpose:** Framework error types.

**Files:**
- `errors.go` - Error type definitions

---

#### `pkg/utils/` - Utilities

**Purpose:** Utility functions.

**Files:**
- `strcase.go` - String case conversion
- `uuid.go` - UUID utilities

---

#### `pkg/registry/` - Extension Registry

**Purpose:** Plugin and extension registry.

**Files:**
- `registry.go` - Registry implementation
- `plugin.go` - Plugin interface
- `extensions.go` - Extension management

---

#### `pkg/models/` - Base Models

**Purpose:** Base model definitions.

**Files:**
- `user.go` - Base user model
- `base_user.go` - Base user schema
- `*.sql` - SQL definitions

---

#### `pkg/cli/` - Command-Line Interface

**Purpose:** CLI tools for framework operations.

**Files:**
- `root.go` - Root command setup
- `registry.go` - Command registry

**Subdirectories:**
- `cmd/` - Command definitions
- `commands/` - Command implementations
  - `project/` - Project commands
    - `new.go` - Create new project
    - `add_app.go` - Add new app
    - `add_model.go` - Add new model
    - `add_api.go` - Add API endpoint
    - `add_handler.go` - Add HTTP handler
    - `add_service.go` - Add service
    - `add_group.go` - Add group
    - `auth.go` - Authentication setup
  - `generation/` - Code generation
    - `generate.go` - Generate code command
  - `migrations/` - Migration commands
    - `makemigrations.go` - Create migrations
    - `up.go` - Apply migrations
    - `rollback.go` - Rollback migrations
    - `status.go` - Migration status
    - `show.go` - Show migration
    - `lint.go` - Lint migrations
    - `squash.go` - Squash migrations
    - `fake.go` - Fake migration
    - `force.go` - Force migration
    - `group.go` - Group migrations
    - `helpers.go` - Migration helpers
  - `server/` - Server commands
    - `runserver.go` - Run development server
  - `development/` - Development tools
    - `shell.go` - Interactive shell
    - `test.go` - Test runner
  - `admin/` - Admin commands
    - `createsuperuser.go` - Create admin user
- `templates/` - Code templates
  - `*.tmpl` - Template files for project generation
- `internal/` - Internal CLI utilities

---

## Command Structure (`cmd/`)

### `cli/cmd/main.go`

**Purpose:** Main entry point for the `forge` CLI tool.

**Functionality:**
- Initializes CLI root command
- Registers all subcommands

**Location:** `forge/cli/cmd/main.go`
- Handles command execution
- Error handling and exit codes

---

## Examples (`examples/`)

### `examples/ecommerce/`

**Purpose:** Complete e-commerce example application.

**Structure:**
- `models/` - E-commerce models (Product, Order, Customer, etc.)
- `api/` - REST API implementation
- `admin/` - Admin configuration
- `cmd/server/` - Server entry point
- `migrations/` - Database migrations
- `config/` - Configuration files

**Files:**
- `README.md` - Example documentation
- `API_EXAMPLES.md` - API usage examples
- `SETUP.md` - Setup instructions

### `examples/library/`

**Purpose:** Library management example.

**Structure:**
- `models/` - Library models (Book, Author, Loan, etc.)
- `migrations/` - Database migrations
- `config/` - Configuration
- `query_examples.go` - Query examples
- `QUERY_EXAMPLES_README.md` - Query examples documentation

---

## Tests (`tests/`)

**Purpose:** Test packages for framework components.

**Structure:**
- `pkg_migrations/` - Migration system tests
- `pkg_query/` - Query system tests
- `pkg_schema/` - Schema system tests
- `cmd_forge/` - CLI tests
- `testhelpers/` - Test helper utilities
- `testdata/` - Test data files

---

## Documentation (`docs/`)

**Purpose:** Framework documentation.

**Key Files:**
- `INDEX.md` - Documentation index
- `ARCHITECTURE.md` - Framework architecture
- `FEATURES.md` - Feature documentation
- `ROADMAP.md` - Development roadmap
- `GETTING_STARTED.md` - Quick start guide
- `USAGE_GUIDE.md` - Usage tutorials
- `API_REFERENCE.md` - API reference
- `SCHEMA_REFERENCE.md` - Schema reference
- `pkg_admin_README.md` - Admin system docs
- `pkg_users_README.md` - User system docs
- And more...

---

## File Naming Conventions

### Go Files
- `*_test.go` - Test files
- `*_v2.go` - Version 2 implementations
- `*_enhanced.go` - Enhanced versions
- `*_helpers.go` - Helper functions
- `*_example.go` - Example code
- `*_complete.go` - Complete examples

### Template Files
- `*.tmpl` - Code generation templates

### SQL Files
- `*.up.sql` - Migration up scripts
- `*.down.sql` - Migration down scripts

---

## Package Dependencies

### Core Dependencies Flow

```
schema → generator → query → db
  ↓         ↓          ↓
admin    models    migrations
  ↓
  api
  ↓
users
```

### Key Dependencies

- `pkg/query` depends on: `pkg/schema`, `pkg/db`
- `pkg/admin` depends on: `pkg/query`, `pkg/schema`
- `pkg/api` depends on: `pkg/query`, `pkg/http`
- `pkg/users` depends on: `pkg/query`, `pkg/schema`
- `pkg/generator` depends on: `pkg/schema`
- `pkg/migrations` depends on: `pkg/schema`, `pkg/db`

---

## Code Organization Principles

1. **Single Responsibility** - Each package has one clear purpose
2. **Dependency Direction** - Lower-level packages don't depend on higher-level ones
3. **Interface-Based** - Packages expose interfaces, not implementations
4. **Type Safety** - Generics used throughout for type safety
5. **Testability** - Packages designed for easy testing
6. **Extensibility** - Everything can be extended or overridden

---

**Last Updated:** January 2025
