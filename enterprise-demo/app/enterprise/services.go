package enterprise

import (
	"context"

	"enterprise-demo/app/enterprise/gen"
)

// ============================================================================
// Service Layer - Business Logic with Complex Queries
// ============================================================================

// OrganizationService provides business logic for organizations
type OrganizationService struct {
	orgRepo *OrganizationRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(orgRepo *OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		orgRepo: orgRepo,
	}
}

// GetOrganizationStats retrieves statistics for an organization
func (s *OrganizationService) GetOrganizationStats(ctx context.Context, orgID int64) (*OrganizationStats, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Get employee count
	empRepo, err := NewEmployeeRepository(s.orgRepo.db)
	if err != nil {
		return nil, err
	}

	employees, err := empRepo.GetByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Get project count
	projRepo, err := NewProjectRepository(s.orgRepo.db)
	if err != nil {
		return nil, err
	}

	projects, err := projRepo.GetByOrganization(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &OrganizationStats{
		Organization:    org,
		EmployeeCount:   int64(len(employees)),
		ProjectCount:    int64(len(projects)),
		ActiveEmployees: s.countActiveEmployees(employees),
		ActiveProjects:  s.countActiveProjects(projects),
	}, nil
}

func (s *OrganizationService) countActiveEmployees(employees []*gen.Employee) int64 {
	count := int64(0)
	for range employees {
		// Access via reflection or direct field access
		// For now, simplified
		count++
	}
	return count
}

func (s *OrganizationService) countActiveProjects(projects []*gen.Project) int64 {
	count := int64(0)
	for range projects {
		// Access via reflection or direct field access
		count++
	}
	return count
}

// OrganizationStats contains organization statistics
type OrganizationStats struct {
	Organization    *gen.Organization
	EmployeeCount   int64
	ProjectCount    int64
	ActiveEmployees int64
	ActiveProjects  int64
}

// ============================================================================
// Employee Service
// ============================================================================

// EmployeeService provides business logic for employees
type EmployeeService struct {
	empRepo *EmployeeRepository
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(empRepo *EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		empRepo: empRepo,
	}
}

// GetEmployeeWorkload retrieves an employee's workload (tasks)
func (s *EmployeeService) GetEmployeeWorkload(ctx context.Context, empID int64) (*EmployeeWorkload, error) {
	taskRepo, err := NewTaskRepository(s.empRepo.db)
	if err != nil {
		return nil, err
	}

	tasks, err := taskRepo.GetByAssignee(ctx, empID)
	if err != nil {
		return nil, err
	}

	pendingTasks := make([]*gen.Task, 0)
	completedTasks := make([]*gen.Task, 0)

	for _, task := range tasks {
		// Check if completed (would use reflection or field access)
		// For now, simplified
		pendingTasks = append(pendingTasks, task)
	}

	return &EmployeeWorkload{
		TotalTasks:     int64(len(tasks)),
		PendingTasks:   int64(len(pendingTasks)),
		CompletedTasks: int64(len(completedTasks)),
		Tasks:          tasks,
	}, nil
}

// EmployeeWorkload contains employee workload information
type EmployeeWorkload struct {
	TotalTasks     int64
	PendingTasks   int64
	CompletedTasks int64
	Tasks          []*gen.Task
}

// ============================================================================
// Project Service
// ============================================================================

// ProjectService provides business logic for projects
type ProjectService struct {
	projRepo *ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(projRepo *ProjectRepository) *ProjectService {
	return &ProjectService{
		projRepo: projRepo,
	}
}

// GetProjectProgress retrieves project progress statistics
func (s *ProjectService) GetProjectProgress(ctx context.Context, projectID int64) (*ProjectProgress, error) {
	project, err := s.projRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	taskRepo, err := NewTaskRepository(s.projRepo.db)
	if err != nil {
		return nil, err
	}

	allTasks, err := taskRepo.GetByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	completedTasks, err := taskRepo.GetCompletedTasks(ctx, projectID)
	if err != nil {
		return nil, err
	}

	totalTasks := int64(len(allTasks))
	completedCount := int64(len(completedTasks))

	var progressPercent float64
	if totalTasks > 0 {
		progressPercent = float64(completedCount) / float64(totalTasks) * 100
	}

	return &ProjectProgress{
		Project:         project,
		TotalTasks:      totalTasks,
		CompletedTasks:  completedCount,
		PendingTasks:    totalTasks - completedCount,
		ProgressPercent: progressPercent,
	}, nil
}

// ProjectProgress contains project progress information
type ProjectProgress struct {
	Project         *gen.Project
	TotalTasks      int64
	CompletedTasks  int64
	PendingTasks    int64
	ProgressPercent float64
}
