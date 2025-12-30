package enterprise

import (
	"context"
	"fmt"

	"enterprise-demo/app/enterprise/gen"
	"github.com/forgego/forge/db"
	query "github.com/forgego/forge/orm"
)

// ============================================================================
// Repository Pattern - Provides data access layer with complex queries
// ============================================================================

// OrganizationRepository provides data access for Organization model
type OrganizationRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Organization]
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(database *db.DB) (*OrganizationRepository, error) {
	manager, err := query.NewManager[gen.Organization]("organizations")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &OrganizationRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves an organization by ID
func (r *OrganizationRepository) GetByID(ctx context.Context, id int64) (*gen.Organization, error) {
	return r.manager.Get(ctx, id)
}

// GetBySlug retrieves an organization by slug
func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*gen.Organization, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	slugField := query.FieldFor[gen.Organization, string](fa, "slug")
	qs, err := r.manager.Filter(slugField.Eq(slug))
	if err != nil {
		return nil, err
	}

	return qs.Get(ctx)
}

// GetActiveOrganizations retrieves all active organizations
func (r *OrganizationRepository) GetActiveOrganizations(ctx context.Context) ([]*gen.Organization, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")
	qs, err := r.manager.Filter(isActiveField.Eq(true))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("created_at")).All(ctx)
}

// GetOrganizationsByTier retrieves organizations by subscription tier
func (r *OrganizationRepository) GetOrganizationsByTier(ctx context.Context, tierID int64) ([]*gen.Organization, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	tierField := query.FieldFor[gen.Organization, int64](fa, "subscription_tier_id")
	qs, err := r.manager.Filter(tierField.Eq(tierID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("created_at")).All(ctx)
}

// SearchOrganizations searches organizations by name or domain
func (r *OrganizationRepository) SearchOrganizations(ctx context.Context, searchTerm string) ([]*gen.Organization, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	nameField := query.FieldFor[gen.Organization, string](fa, "name")
	domainField := query.FieldFor[gen.Organization, string](fa, "domain")

	// Complex OR query using Q objects
	searchQ := query.NewQ(nameField.Contains(searchTerm)).
		Or(query.NewQ(domainField.Contains(searchTerm)))

	qs, err := r.manager.Filter(searchQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("created_at")).All(ctx)
}

// GetOrganizationsWithStats retrieves organizations with statistics
func (r *OrganizationRepository) GetOrganizationsWithStats(ctx context.Context, limit, offset int) ([]*gen.Organization, int64, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, 0, err
	}

	isActiveField := query.FieldFor[gen.Organization, bool](fa, "is_active")
	qs, err := r.manager.Filter(isActiveField.Eq(true))
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := qs.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	results, err := qs.
		OrderBy(query.Desc("created_at")).
		Limit(limit).
		Offset(offset).
		All(ctx)

	return results, total, err
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *gen.Organization) error {
	return r.manager.Create(ctx, org)
}

// Update updates an organization
func (r *OrganizationRepository) Update(ctx context.Context, org *gen.Organization) error {
	return r.manager.Update(ctx, org)
}

// Delete deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, org *gen.Organization) error {
	return r.manager.Delete(ctx, org)
}

// ============================================================================
// Employee Repository
// ============================================================================

// EmployeeRepository provides data access for Employee model
type EmployeeRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Employee]
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(database *db.DB) (*EmployeeRepository, error) {
	manager, err := query.NewManager[gen.Employee]("employees")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &EmployeeRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves an employee by ID
func (r *EmployeeRepository) GetByID(ctx context.Context, id int64) (*gen.Employee, error) {
	return r.manager.Get(ctx, id)
}

// GetByEmail retrieves an employee by email
func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	emailField := query.FieldFor[gen.Employee, string](fa, "email")
	qs, err := r.manager.Filter(emailField.Eq(email))
	if err != nil {
		return nil, err
	}

	return qs.Get(ctx)
}

// GetByOrganization retrieves all employees for an organization
func (r *EmployeeRepository) GetByOrganization(ctx context.Context, orgID int64) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	qs, err := r.manager.Filter(orgField.Eq(orgID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("last_name"), query.Asc("first_name")).All(ctx)
}

// GetByDepartment retrieves employees in a department
func (r *EmployeeRepository) GetByDepartment(ctx context.Context, deptID int64) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	deptField := query.FieldFor[gen.Employee, int64](fa, "department_id")
	qs, err := r.manager.Filter(deptField.Eq(deptID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("last_name"), query.Asc("first_name")).All(ctx)
}

// GetActiveEmployees retrieves all active employees
func (r *EmployeeRepository) GetActiveEmployees(ctx context.Context, orgID int64) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	isActiveField := query.FieldFor[gen.Employee, bool](fa, "is_active")

	// Complex AND query
	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(isActiveField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("last_name"), query.Asc("first_name")).All(ctx)
}

// GetManagers retrieves all managers
func (r *EmployeeRepository) GetManagers(ctx context.Context, orgID int64) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	isManagerField := query.FieldFor[gen.Employee, bool](fa, "is_manager")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(isManagerField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("last_name"), query.Asc("first_name")).All(ctx)
}

// SearchEmployees searches employees by name or email
func (r *EmployeeRepository) SearchEmployees(ctx context.Context, orgID int64, searchTerm string) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	firstNameField := query.FieldFor[gen.Employee, string](fa, "first_name")
	lastNameField := query.FieldFor[gen.Employee, string](fa, "last_name")
	emailField := query.FieldFor[gen.Employee, string](fa, "email")

	// Complex query: (orgID matches) AND (name or email contains search term)
	orgQ := query.NewQ(orgField.Eq(orgID))
	searchQ := query.NewQ(firstNameField.Contains(searchTerm)).
		Or(query.NewQ(lastNameField.Contains(searchTerm))).
		Or(query.NewQ(emailField.Contains(searchTerm)))

	filterQ := orgQ.And(searchQ)

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("last_name"), query.Asc("first_name")).All(ctx)
}

// GetEmployeesBySalaryRange retrieves employees within a salary range
func (r *EmployeeRepository) GetEmployeesBySalaryRange(ctx context.Context, orgID int64, minSalary, maxSalary float64) ([]*gen.Employee, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	salaryField := query.FieldFor[gen.Employee, float64](fa, "salary")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(salaryField.Gte(minSalary))).
		And(query.NewQ(salaryField.Lte(maxSalary)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("salary")).All(ctx)
}

// GetEmployeesWithPagination retrieves employees with pagination
func (r *EmployeeRepository) GetEmployeesWithPagination(ctx context.Context, orgID int64, limit, offset int) ([]*gen.Employee, int64, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, 0, err
	}

	orgField := query.FieldFor[gen.Employee, int64](fa, "organization_id")
	qs, err := r.manager.Filter(orgField.Eq(orgID))
	if err != nil {
		return nil, 0, err
	}

	total, err := qs.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	results, err := qs.
		OrderBy(query.Asc("last_name"), query.Asc("first_name")).
		Limit(limit).
		Offset(offset).
		All(ctx)

	return results, total, err
}

// Create creates a new employee
func (r *EmployeeRepository) Create(ctx context.Context, emp *gen.Employee) error {
	return r.manager.Create(ctx, emp)
}

// Update updates an employee
func (r *EmployeeRepository) Update(ctx context.Context, emp *gen.Employee) error {
	return r.manager.Update(ctx, emp)
}

// Delete deletes an employee
func (r *EmployeeRepository) Delete(ctx context.Context, emp *gen.Employee) error {
	return r.manager.Delete(ctx, emp)
}

// ============================================================================
// Project Repository
// ============================================================================

// ProjectRepository provides data access for Project model
type ProjectRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Project]
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(database *db.DB) (*ProjectRepository, error) {
	manager, err := query.NewManager[gen.Project]("projects")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &ProjectRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*gen.Project, error) {
	return r.manager.Get(ctx, id)
}

// GetByOrganization retrieves all projects for an organization
func (r *ProjectRepository) GetByOrganization(ctx context.Context, orgID int64) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Project, int64](fa, "organization_id")
	qs, err := r.manager.Filter(orgField.Eq(orgID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Desc("created_at")).All(ctx)
}

// GetActiveProjects retrieves active projects
func (r *ProjectRepository) GetActiveProjects(ctx context.Context, orgID int64) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Project, int64](fa, "organization_id")
	isActiveField := query.FieldFor[gen.Project, bool](fa, "is_active")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(isActiveField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Desc("created_at")).All(ctx)
}

// GetProjectsByStatus retrieves projects by status
func (r *ProjectRepository) GetProjectsByStatus(ctx context.Context, orgID int64, status string) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Project, int64](fa, "organization_id")
	statusField := query.FieldFor[gen.Project, string](fa, "status")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(statusField.Eq(status)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Desc("created_at")).All(ctx)
}

// GetProjectsByManager retrieves projects managed by an employee
func (r *ProjectRepository) GetProjectsByManager(ctx context.Context, managerID int64) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	managerField := query.FieldFor[gen.Project, int64](fa, "project_manager_id")
	qs, err := r.manager.Filter(managerField.Eq(managerID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Desc("created_at")).All(ctx)
}

// GetProjectsByBudgetRange retrieves projects within a budget range
func (r *ProjectRepository) GetProjectsByBudgetRange(ctx context.Context, orgID int64, minBudget, maxBudget float64) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Project, int64](fa, "organization_id")
	budgetField := query.FieldFor[gen.Project, float64](fa, "budget")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(budgetField.Gte(minBudget))).
		And(query.NewQ(budgetField.Lte(maxBudget)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("budget")).All(ctx)
}

// SearchProjects searches projects by name or code
func (r *ProjectRepository) SearchProjects(ctx context.Context, orgID int64, searchTerm string) ([]*gen.Project, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Project, int64](fa, "organization_id")
	nameField := query.FieldFor[gen.Project, string](fa, "name")
	codeField := query.FieldFor[gen.Project, string](fa, "code")

	orgQ := query.NewQ(orgField.Eq(orgID))
	searchQ := query.NewQ(nameField.Contains(searchTerm)).
		Or(query.NewQ(codeField.Contains(searchTerm)))

	filterQ := orgQ.And(searchQ)

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Desc("created_at")).All(ctx)
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, proj *gen.Project) error {
	return r.manager.Create(ctx, proj)
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, proj *gen.Project) error {
	return r.manager.Update(ctx, proj)
}

// Delete deletes a project
func (r *ProjectRepository) Delete(ctx context.Context, proj *gen.Project) error {
	return r.manager.Delete(ctx, proj)
}

// ============================================================================
// Task Repository
// ============================================================================

// TaskRepository provides data access for Task model
type TaskRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Task]
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(database *db.DB) (*TaskRepository, error) {
	manager, err := query.NewManager[gen.Task]("tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &TaskRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*gen.Task, error) {
	return r.manager.Get(ctx, id)
}

// GetByProject retrieves all tasks for a project
func (r *TaskRepository) GetByProject(ctx context.Context, projectID int64) ([]*gen.Task, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	projectField := query.FieldFor[gen.Task, int64](fa, "project_id")
	qs, err := r.manager.Filter(projectField.Eq(projectID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("priority"), query.Asc("due_date")).All(ctx)
}

// GetByAssignee retrieves tasks assigned to an employee
func (r *TaskRepository) GetByAssignee(ctx context.Context, employeeID int64) ([]*gen.Task, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	assigneeField := query.FieldFor[gen.Task, int64](fa, "assigned_to_id")
	qs, err := r.manager.Filter(assigneeField.Eq(employeeID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("due_date"), query.Desc("priority")).All(ctx)
}

// GetCompletedTasks retrieves completed tasks
func (r *TaskRepository) GetCompletedTasks(ctx context.Context, projectID int64) ([]*gen.Task, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	projectField := query.FieldFor[gen.Task, int64](fa, "project_id")
	isCompletedField := query.FieldFor[gen.Task, bool](fa, "is_completed")

	filterQ := query.NewQ(projectField.Eq(projectID)).
		And(query.NewQ(isCompletedField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Desc("completed_at")).All(ctx)
}

// GetPendingTasks retrieves pending (not completed) tasks
func (r *TaskRepository) GetPendingTasks(ctx context.Context, projectID int64) ([]*gen.Task, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	projectField := query.FieldFor[gen.Task, int64](fa, "project_id")
	isCompletedField := query.FieldFor[gen.Task, bool](fa, "is_completed")

	filterQ := query.NewQ(projectField.Eq(projectID)).
		And(query.NewQ(isCompletedField.Eq(false)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("due_date"), query.Desc("priority")).All(ctx)
}

// GetOverdueTasks retrieves overdue tasks
func (r *TaskRepository) GetOverdueTasks(ctx context.Context, projectID int64) ([]*gen.Task, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	projectField := query.FieldFor[gen.Task, int64](fa, "project_id")
	isCompletedField := query.FieldFor[gen.Task, bool](fa, "is_completed")
	// Note: dueDateField would need time.Time in real implementation
	// dueDateField := query.FieldFor[gen.Task, time.Time](fa, "due_date")

	// This is a simplified example - would need proper date comparison
	filterQ := query.NewQ(projectField.Eq(projectID)).
		And(query.NewQ(isCompletedField.Eq(false)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("due_date")).All(ctx)
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *gen.Task) error {
	return r.manager.Create(ctx, task)
}

// Update updates a task
func (r *TaskRepository) Update(ctx context.Context, task *gen.Task) error {
	return r.manager.Update(ctx, task)
}

// Delete deletes a task
func (r *TaskRepository) Delete(ctx context.Context, task *gen.Task) error {
	return r.manager.Delete(ctx, task)
}
