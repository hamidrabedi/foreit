-- Create schema_migrations table for tracking applied migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    checksum TEXT,
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Create table: organizations
CREATE TABLE IF NOT EXISTS organizations (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(100) NOT NULL UNIQUE,
    "domain" VARCHAR(255) UNIQUE,
    "description" TEXT,
    "industry" VARCHAR(100),
    "website" VARCHAR(500),
    "logo_url" VARCHAR(500),
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_verified" BOOLEAN DEFAULT FALSE,
    "subscription_tier_id" BIGINT,
    "trial_ends_at" TIME,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: subscription_tiers
CREATE TABLE IF NOT EXISTS subscription_tiers (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(100) NOT NULL,
    "slug" VARCHAR(50) NOT NULL UNIQUE,
    "description" TEXT,
    "monthly_price" DOUBLE PRECISION DEFAULT 0.000000,
    "yearly_price" DOUBLE PRECISION DEFAULT 0.000000,
    "max_users" BIGINT DEFAULT 10,
    "max_projects" BIGINT DEFAULT 5,
    "max_storage_gb" BIGINT DEFAULT 10,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_featured" BOOLEAN DEFAULT FALSE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: departments
CREATE TABLE IF NOT EXISTS departments (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "organization_id" BIGINT NOT NULL,
    "name" VARCHAR(200) NOT NULL,
    "code" VARCHAR(50),
    "description" TEXT,
    "parent_department_id" BIGINT,
    "manager_id" BIGINT,
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: employees
CREATE TABLE IF NOT EXISTS employees (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "organization_id" BIGINT NOT NULL,
    "department_id" BIGINT,
    "employee_number" VARCHAR(50) UNIQUE,
    "first_name" VARCHAR(100) NOT NULL,
    "last_name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "phone" VARCHAR(20),
    "job_title" VARCHAR(200),
    "employee_type" VARCHAR(50) DEFAULT 'full_time',
    "salary" DOUBLE PRECISION DEFAULT 0.000000,
    "hire_date" TIME,
    "termination_date" TIME,
    "is_active" BOOLEAN DEFAULT TRUE,
    "is_manager" BOOLEAN DEFAULT FALSE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: projects
CREATE TABLE IF NOT EXISTS projects (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "organization_id" BIGINT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "code" VARCHAR(50),
    "description" TEXT,
    "status" VARCHAR(50) DEFAULT 'planning',
    "project_manager_id" BIGINT,
    "client_id" BIGINT,
    "budget" DOUBLE PRECISION DEFAULT 0.000000,
    "start_date" TIME,
    "end_date" TIME,
    "deadline" TIME,
    "priority" BIGINT DEFAULT 0,
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: clients
CREATE TABLE IF NOT EXISTS clients (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "organization_id" BIGINT NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "company_name" VARCHAR(255),
    "email" VARCHAR(255),
    "phone" VARCHAR(20),
    "address" VARCHAR(500),
    "city" VARCHAR(100),
    "state" VARCHAR(100),
    "zip_code" VARCHAR(20),
    "country" VARCHAR(100),
    "industry" VARCHAR(100),
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: tasks
CREATE TABLE IF NOT EXISTS tasks (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "project_id" BIGINT NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "status" VARCHAR(50) DEFAULT 'todo',
    "priority" VARCHAR(20) DEFAULT 'medium',
    "assigned_to_id" BIGINT,
    "created_by_id" BIGINT,
    "due_date" TIME,
    "started_at" TIME,
    "completed_at" TIME,
    "estimated_hours" BIGINT DEFAULT 0,
    "actual_hours" BIGINT DEFAULT 0,
    "is_completed" BOOLEAN DEFAULT FALSE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

-- Create table: skills
CREATE TABLE IF NOT EXISTS skills (
    "id" BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "name" VARCHAR(100) NOT NULL UNIQUE,
    "category" VARCHAR(100),
    "description" TEXT,
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIME DEFAULT now() NOT NULL,
    "updated_at" TIME DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_org_slug ON organizations ("slug");

CREATE UNIQUE INDEX IF NOT EXISTS idx_org_domain ON organizations ("domain");

CREATE INDEX IF NOT EXISTS idx_org_subscription ON organizations ("subscription_tier_id");

CREATE INDEX IF NOT EXISTS idx_org_active ON organizations ("is_active");

CREATE UNIQUE INDEX IF NOT EXISTS idx_tier_slug ON subscription_tiers ("slug");

CREATE INDEX IF NOT EXISTS idx_tier_active ON subscription_tiers ("is_active");

CREATE INDEX IF NOT EXISTS idx_dept_org ON departments ("organization_id");

CREATE INDEX IF NOT EXISTS idx_dept_parent ON departments ("parent_department_id");

CREATE INDEX IF NOT EXISTS idx_dept_manager ON departments ("manager_id");

CREATE UNIQUE INDEX IF NOT EXISTS idx_dept_code ON departments ("organization_id", "code");

CREATE INDEX IF NOT EXISTS idx_emp_org ON employees ("organization_id");

CREATE INDEX IF NOT EXISTS idx_emp_dept ON employees ("department_id");

CREATE UNIQUE INDEX IF NOT EXISTS idx_emp_email ON employees ("email");

CREATE UNIQUE INDEX IF NOT EXISTS idx_emp_number ON employees ("employee_number");

CREATE INDEX IF NOT EXISTS idx_emp_active ON employees ("is_active");

CREATE INDEX IF NOT EXISTS idx_proj_org ON projects ("organization_id");

CREATE INDEX IF NOT EXISTS idx_proj_manager ON projects ("project_manager_id");

CREATE INDEX IF NOT EXISTS idx_proj_client ON projects ("client_id");

CREATE INDEX IF NOT EXISTS idx_proj_status ON projects ("status");

CREATE UNIQUE INDEX IF NOT EXISTS idx_proj_code ON projects ("organization_id", "code");

CREATE INDEX IF NOT EXISTS idx_client_org ON clients ("organization_id");

CREATE INDEX IF NOT EXISTS idx_client_email ON clients ("email");

CREATE INDEX IF NOT EXISTS idx_task_project ON tasks ("project_id");

CREATE INDEX IF NOT EXISTS idx_task_assigned ON tasks ("assigned_to_id");

CREATE INDEX IF NOT EXISTS idx_task_status ON tasks ("status");

CREATE INDEX IF NOT EXISTS idx_task_completed ON tasks ("is_completed");

CREATE UNIQUE INDEX IF NOT EXISTS idx_skill_name ON skills ("name");

CREATE INDEX IF NOT EXISTS idx_skill_category ON skills ("category");
