package enterprise

import (
	"github.com/forgego/forge/schema"
)

// ============================================================================
// Organization Model - Represents a company/organization
// ============================================================================

type Organization struct {
	schema.BaseSchema
}

func (Organization) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(255).VerboseName("Organization Name").Build(),
		schema.String("slug").Required().MaxLength(100).Unique().VerboseName("URL Slug").Build(),
		schema.String("domain").MaxLength(255).Unique().VerboseName("Domain").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.String("industry").MaxLength(100).VerboseName("Industry").Build(),
		schema.String("website").MaxLength(500).VerboseName("Website").Build(),
		schema.String("logo_url").MaxLength(500).VerboseName("Logo URL").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Bool("is_verified").Default(false).VerboseName("Verified").Build(),
		schema.Int64("subscription_tier_id").VerboseName("Subscription Tier ID").Build(),
		schema.Time("trial_ends_at").VerboseName("Trial Ends At").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Organization) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "organizations",
		VerboseName:       "Organization",
		VerboseNamePlural: "Organizations",
		OrderBy:           []string{"-created_at", "name"},
		Indexes: []schema.Index{
			{Name: "idx_org_slug", Fields: []string{"slug"}, Unique: true},
			{Name: "idx_org_domain", Fields: []string{"domain"}, Unique: true},
			{Name: "idx_org_subscription", Fields: []string{"subscription_tier_id"}, Unique: false},
			{Name: "idx_org_active", Fields: []string{"is_active"}, Unique: false},
		},
	}
}

func (Organization) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("subscription_tier_id", "SubscriptionTier").
			OnDelete(schema.CascadeCASCADE).
			Build(),
	}
}

func (Organization) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// SubscriptionTier Model - Represents subscription plans
// ============================================================================

type SubscriptionTier struct {
	schema.BaseSchema
}

func (SubscriptionTier) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(100).VerboseName("Tier Name").Build(),
		schema.String("slug").Required().MaxLength(50).Unique().VerboseName("Slug").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.Float64("monthly_price").Default(0.0).VerboseName("Monthly Price").Build(),
		schema.Float64("yearly_price").Default(0.0).VerboseName("Yearly Price").Build(),
		schema.Int64("max_users").Default(10).VerboseName("Max Users").Build(),
		schema.Int64("max_projects").Default(5).VerboseName("Max Projects").Build(),
		schema.Int64("max_storage_gb").Default(10).VerboseName("Max Storage (GB)").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Bool("is_featured").Default(false).VerboseName("Featured").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (SubscriptionTier) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "subscription_tiers",
		VerboseName:       "Subscription Tier",
		VerboseNamePlural: "Subscription Tiers",
		OrderBy:           []string{"monthly_price", "name"},
		Indexes: []schema.Index{
			{Name: "idx_tier_slug", Fields: []string{"slug"}, Unique: true},
			{Name: "idx_tier_active", Fields: []string{"is_active"}, Unique: false},
		},
	}
}

func (SubscriptionTier) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (SubscriptionTier) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Department Model - Represents departments within organizations
// ============================================================================

type Department struct {
	schema.BaseSchema
}

func (Department) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("organization_id").Required().VerboseName("Organization ID").Build(),
		schema.String("name").Required().MaxLength(200).VerboseName("Department Name").Build(),
		schema.String("code").MaxLength(50).VerboseName("Department Code").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.Int64("parent_department_id").VerboseName("Parent Department ID").Build(),
		schema.Int64("manager_id").VerboseName("Manager ID").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Department) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "departments",
		VerboseName:       "Department",
		VerboseNamePlural: "Departments",
		OrderBy:           []string{"organization_id", "name"},
		Indexes: []schema.Index{
			{Name: "idx_dept_org", Fields: []string{"organization_id"}, Unique: false},
			{Name: "idx_dept_parent", Fields: []string{"parent_department_id"}, Unique: false},
			{Name: "idx_dept_manager", Fields: []string{"manager_id"}, Unique: false},
			{Name: "idx_dept_code", Fields: []string{"organization_id", "code"}, Unique: true},
		},
	}
}

func (Department) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("organization_id", "Organization").
			OnDelete(schema.CascadeCASCADE).
			Build(),
		schema.ForeignKey("parent_department_id", "Department").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
		schema.ForeignKey("manager_id", "Employee").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
	}
}

func (Department) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Employee Model - Represents employees/users
// ============================================================================

type Employee struct {
	schema.BaseSchema
}

func (Employee) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("organization_id").Required().VerboseName("Organization ID").Build(),
		schema.Int64("department_id").VerboseName("Department ID").Build(),
		schema.String("employee_number").MaxLength(50).Unique().VerboseName("Employee Number").Build(),
		schema.String("first_name").Required().MaxLength(100).VerboseName("First Name").Build(),
		schema.String("last_name").Required().MaxLength(100).VerboseName("Last Name").Build(),
		schema.String("email").Required().MaxLength(255).Unique().VerboseName("Email").Build(),
		schema.String("phone").MaxLength(20).VerboseName("Phone").Build(),
		schema.String("job_title").MaxLength(200).VerboseName("Job Title").Build(),
		schema.String("employee_type").MaxLength(50).Default("full_time").VerboseName("Employee Type").Build(),
		schema.Float64("salary").Default(0.0).VerboseName("Salary").Build(),
		schema.Time("hire_date").VerboseName("Hire Date").Build(),
		schema.Time("termination_date").VerboseName("Termination Date").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Bool("is_manager").Default(false).VerboseName("Is Manager").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Employee) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "employees",
		VerboseName:       "Employee",
		VerboseNamePlural: "Employees",
		OrderBy:           []string{"-created_at", "last_name", "first_name"},
		Indexes: []schema.Index{
			{Name: "idx_emp_org", Fields: []string{"organization_id"}, Unique: false},
			{Name: "idx_emp_dept", Fields: []string{"department_id"}, Unique: false},
			{Name: "idx_emp_email", Fields: []string{"email"}, Unique: true},
			{Name: "idx_emp_number", Fields: []string{"employee_number"}, Unique: true},
			{Name: "idx_emp_active", Fields: []string{"is_active"}, Unique: false},
		},
	}
}

func (Employee) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("organization_id", "Organization").
			OnDelete(schema.CascadeCASCADE).
			Build(),
		schema.ForeignKey("department_id", "Department").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
		schema.ManyToMany("projects", "Project").
			ThroughTable("employee_projects").
			Build(),
		schema.ManyToMany("skills", "Skill").
			ThroughTable("employee_skills").
			Build(),
	}
}

func (Employee) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Project Model - Represents projects
// ============================================================================

type Project struct {
	schema.BaseSchema
}

func (Project) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("organization_id").Required().VerboseName("Organization ID").Build(),
		schema.String("name").Required().MaxLength(255).VerboseName("Project Name").Build(),
		schema.String("code").MaxLength(50).VerboseName("Project Code").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.String("status").MaxLength(50).Default("planning").VerboseName("Status").Build(),
		schema.Int64("project_manager_id").VerboseName("Project Manager ID").Build(),
		schema.Int64("client_id").VerboseName("Client ID").Build(),
		schema.Float64("budget").Default(0.0).VerboseName("Budget").Build(),
		schema.Time("start_date").VerboseName("Start Date").Build(),
		schema.Time("end_date").VerboseName("End Date").Build(),
		schema.Time("deadline").VerboseName("Deadline").Build(),
		schema.Int64("priority").Default(0).VerboseName("Priority").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Project) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "projects",
		VerboseName:       "Project",
		VerboseNamePlural: "Projects",
		OrderBy:           []string{"-priority", "-created_at", "name"},
		Indexes: []schema.Index{
			{Name: "idx_proj_org", Fields: []string{"organization_id"}, Unique: false},
			{Name: "idx_proj_manager", Fields: []string{"project_manager_id"}, Unique: false},
			{Name: "idx_proj_client", Fields: []string{"client_id"}, Unique: false},
			{Name: "idx_proj_status", Fields: []string{"status"}, Unique: false},
			{Name: "idx_proj_code", Fields: []string{"organization_id", "code"}, Unique: true},
		},
	}
}

func (Project) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("organization_id", "Organization").
			OnDelete(schema.CascadeCASCADE).
			Build(),
		schema.ForeignKey("project_manager_id", "Employee").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
		schema.ForeignKey("client_id", "Client").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
		schema.ManyToMany("employees", "Employee").
			ThroughTable("employee_projects").
			Build(),
	}
}

func (Project) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Client Model - Represents external clients
// ============================================================================

type Client struct {
	schema.BaseSchema
}

func (Client) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("organization_id").Required().VerboseName("Organization ID").Build(),
		schema.String("name").Required().MaxLength(255).VerboseName("Client Name").Build(),
		schema.String("company_name").MaxLength(255).VerboseName("Company Name").Build(),
		schema.String("email").MaxLength(255).VerboseName("Email").Build(),
		schema.String("phone").MaxLength(20).VerboseName("Phone").Build(),
		schema.String("address").MaxLength(500).VerboseName("Address").Build(),
		schema.String("city").MaxLength(100).VerboseName("City").Build(),
		schema.String("state").MaxLength(100).VerboseName("State").Build(),
		schema.String("zip_code").MaxLength(20).VerboseName("Zip Code").Build(),
		schema.String("country").MaxLength(100).VerboseName("Country").Build(),
		schema.String("industry").MaxLength(100).VerboseName("Industry").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Client) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "clients",
		VerboseName:       "Client",
		VerboseNamePlural: "Clients",
		OrderBy:           []string{"name"},
		Indexes: []schema.Index{
			{Name: "idx_client_org", Fields: []string{"organization_id"}, Unique: false},
			{Name: "idx_client_email", Fields: []string{"email"}, Unique: false},
		},
	}
}

func (Client) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("organization_id", "Organization").
			OnDelete(schema.CascadeCASCADE).
			Build(),
	}
}

func (Client) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Task Model - Represents tasks within projects
// ============================================================================

type Task struct {
	schema.BaseSchema
}

func (Task) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("project_id").Required().VerboseName("Project ID").Build(),
		schema.String("title").Required().MaxLength(255).VerboseName("Task Title").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.String("status").MaxLength(50).Default("todo").VerboseName("Status").Build(),
		schema.String("priority").MaxLength(20).Default("medium").VerboseName("Priority").Build(),
		schema.Int64("assigned_to_id").VerboseName("Assigned To ID").Build(),
		schema.Int64("created_by_id").VerboseName("Created By ID").Build(),
		schema.Time("due_date").VerboseName("Due Date").Build(),
		schema.Time("started_at").VerboseName("Started At").Build(),
		schema.Time("completed_at").VerboseName("Completed At").Build(),
		schema.Int64("estimated_hours").Default(0).VerboseName("Estimated Hours").Build(),
		schema.Int64("actual_hours").Default(0).VerboseName("Actual Hours").Build(),
		schema.Bool("is_completed").Default(false).VerboseName("Completed").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Task) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "tasks",
		VerboseName:       "Task",
		VerboseNamePlural: "Tasks",
		OrderBy:           []string{"-priority", "due_date", "title"},
		Indexes: []schema.Index{
			{Name: "idx_task_project", Fields: []string{"project_id"}, Unique: false},
			{Name: "idx_task_assigned", Fields: []string{"assigned_to_id"}, Unique: false},
			{Name: "idx_task_status", Fields: []string{"status"}, Unique: false},
			{Name: "idx_task_completed", Fields: []string{"is_completed"}, Unique: false},
		},
	}
}

func (Task) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("project_id", "Project").
			OnDelete(schema.CascadeCASCADE).
			Build(),
		schema.ForeignKey("assigned_to_id", "Employee").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
		schema.ForeignKey("created_by_id", "Employee").
			OnDelete(schema.CascadeSET_NULL).
			Build(),
	}
}

func (Task) Hooks() *schema.ModelHooks {
	return nil
}

// ============================================================================
// Skill Model - Represents skills (many-to-many with employees)
// ============================================================================

type Skill struct {
	schema.BaseSchema
}

func (Skill) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(100).Unique().VerboseName("Skill Name").Build(),
		schema.String("category").MaxLength(100).VerboseName("Category").Build(),
		schema.Text("description").VerboseName("Description").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Skill) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "skills",
		VerboseName:       "Skill",
		VerboseNamePlural: "Skills",
		OrderBy:           []string{"category", "name"},
		Indexes: []schema.Index{
			{Name: "idx_skill_name", Fields: []string{"name"}, Unique: true},
			{Name: "idx_skill_category", Fields: []string{"category"}, Unique: false},
		},
	}
}

func (Skill) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ManyToMany("employees", "Employee").
			ThroughTable("employee_skills").
			Build(),
	}
}

func (Skill) Hooks() *schema.ModelHooks {
	return nil
}
