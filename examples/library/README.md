# Library Management System

A comprehensive example project demonstrating Forge framework capabilities with:

- Full ORM models with field builders
- SQLite3 database
- Admin interface generation
- Migrations

## Models

- **Author**: Book authors with name, email, bio
- **Book**: Books with title, ISBN, author, category
- **Category**: Book categories
- **Borrower**: Library members
- **Loan**: Book borrowing records

## Getting Started

### Prerequisites

**Note**: SQLite3 requires CGO to be enabled. Make sure you have a C compiler installed:

- Windows: Install MinGW or use TDM-GCC
- Linux: Install `gcc` via package manager
- macOS: Install Xcode Command Line Tools

Build the forge CLI with CGO enabled:

```bash
cd forge
CGO_ENABLED=1 go build -o forge.exe ./cli/cmd
```

### Setup Steps

1. Generate ORM code from models:

   ```bash
   forge generate
   ```

   This generates:

   - Model structs with fields
   - Managers for database operations
   - Querysets for query building
   - Field expressions

2. Create migrations:

   ```bash
   forge makemigrations
   ```

   This creates SQL migration files in `migrations/` directory with proper SQLite syntax.

3. Apply migrations:

   ```bash
   forge migrate
   ```

   This applies all pending migrations to create the database tables.

4. Start the server:

   ```bash
   forge runserver
   ```

5. Access admin interface at: http://localhost:8000/admin/

## Project Structure

- `models/` - Model definitions using schema builders
  - `*.go` - Model definitions
  - `*.gen.go` - Generated ORM code (do not edit)
- `migrations/` - Database migration files
  - `*.up.sql` - Migration up scripts
  - `*.down.sql` - Migration down scripts
- `config/` - Configuration files
  - `config.yaml` - Application configuration
- `main.go` - Application entry point

## Features Demonstrated

1. **Schema Builders**: All models use fluent field builders (Int64, String, Bool, Time, Float64)
2. **Field Options**: Required, MaxLength, Unique, Default, AutoIncrement, Primary
3. **Meta Configuration**: Table names, verbose names, ordering, indexes
4. **SQLite Support**: Proper SQLite syntax in migrations (INTEGER PRIMARY KEY AUTOINCREMENT, etc.)
5. **Generated ORM**: Type-safe models, managers, and querysets
6. **Admin Integration**: Models are automatically available in admin interface
