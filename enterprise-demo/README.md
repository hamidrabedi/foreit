# Enterprise Demo Project

This is a comprehensive enterprise project demonstrating forge's capabilities with complex type-safe schemas, code generation, and migrations.

## Project Structure

```
enterprise-demo/
├── cmd/server/main.go      # Application entry point
├── app/enterprise/
│   └── models.go          # Complex schema definitions
├── config/
│   └── config.yaml        # Configuration
├── migrations/             # Database migrations
└── go.mod
```

## Models

This project includes the following complex models:

1. **Organization** - Companies/organizations with subscriptions
2. **SubscriptionTier** - Subscription plans
3. **Department** - Departments within organizations (self-referential)
4. **Employee** - Employees with many-to-many relationships
5. **Project** - Projects with multiple relationships
6. **Client** - External clients
7. **Task** - Tasks within projects
8. **Skill** - Skills (many-to-many with employees)

## Features Demonstrated

- Complex foreign key relationships
- Self-referential relationships (Department -> Department)
- Many-to-many relationships (Employee <-> Project, Employee <-> Skill)
- Composite unique indexes
- Multiple indexes per model
- Cascade and set null delete behaviors
- Type-safe schema definitions

## Setup

1. **Generate Code:**
   ```bash
   forge generate --models app/enterprise --output app/enterprise
   ```

2. **Create Migrations:**
   ```bash
   forge makemigrations initial --auto --models app/enterprise
   ```

3. **Apply Migrations:**
   ```bash
   forge migrate up
   ```

4. **Run Server:**
   ```bash
   go run cmd/server/main.go
   ```

## Database Schema

The generated migrations will create:
- 8 main tables
- 2 many-to-many junction tables (employee_projects, employee_skills)
- Multiple indexes for performance
- Foreign key constraints with proper cascade behaviors
