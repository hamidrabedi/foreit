package enterprise

import (
	"context"
	"fmt"

	"enterprise-demo/app/enterprise/gen"
	"github.com/forgego/forge/db"
	query "github.com/forgego/forge/orm"
)

// ============================================================================
// Additional Repositories for Remaining Models
// ============================================================================

// DepartmentRepository provides data access for Department model
type DepartmentRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Department]
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository(database *db.DB) (*DepartmentRepository, error) {
	manager, err := query.NewManager[gen.Department]("departments")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &DepartmentRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a department by ID
func (r *DepartmentRepository) GetByID(ctx context.Context, id int64) (*gen.Department, error) {
	return r.manager.Get(ctx, id)
}

// GetByOrganization retrieves all departments for an organization
func (r *DepartmentRepository) GetByOrganization(ctx context.Context, orgID int64) ([]*gen.Department, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Department, int64](fa, "organization_id")
	qs, err := r.manager.Filter(orgField.Eq(orgID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// GetActiveDepartments retrieves active departments
func (r *DepartmentRepository) GetActiveDepartments(ctx context.Context, orgID int64) ([]*gen.Department, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Department, int64](fa, "organization_id")
	isActiveField := query.FieldFor[gen.Department, bool](fa, "is_active")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(isActiveField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// GetByParent retrieves child departments
func (r *DepartmentRepository) GetByParent(ctx context.Context, parentID int64) ([]*gen.Department, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	parentField := query.FieldFor[gen.Department, int64](fa, "parent_department_id")
	qs, err := r.manager.Filter(parentField.Eq(parentID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// Create creates a new department
func (r *DepartmentRepository) Create(ctx context.Context, dept *gen.Department) error {
	return r.manager.Create(ctx, dept)
}

// Update updates a department
func (r *DepartmentRepository) Update(ctx context.Context, dept *gen.Department) error {
	return r.manager.Update(ctx, dept)
}

// Delete deletes a department
func (r *DepartmentRepository) Delete(ctx context.Context, dept *gen.Department) error {
	return r.manager.Delete(ctx, dept)
}

// ============================================================================
// Client Repository
// ============================================================================

// ClientRepository provides data access for Client model
type ClientRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Client]
}

// NewClientRepository creates a new client repository
func NewClientRepository(database *db.DB) (*ClientRepository, error) {
	manager, err := query.NewManager[gen.Client]("clients")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &ClientRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a client by ID
func (r *ClientRepository) GetByID(ctx context.Context, id int64) (*gen.Client, error) {
	return r.manager.Get(ctx, id)
}

// GetByOrganization retrieves all clients for an organization
func (r *ClientRepository) GetByOrganization(ctx context.Context, orgID int64) ([]*gen.Client, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Client, int64](fa, "organization_id")
	qs, err := r.manager.Filter(orgField.Eq(orgID))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// GetActiveClients retrieves active clients
func (r *ClientRepository) GetActiveClients(ctx context.Context, orgID int64) ([]*gen.Client, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Client, int64](fa, "organization_id")
	isActiveField := query.FieldFor[gen.Client, bool](fa, "is_active")

	filterQ := query.NewQ(orgField.Eq(orgID)).
		And(query.NewQ(isActiveField.Eq(true)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// SearchClients searches clients by name or email
func (r *ClientRepository) SearchClients(ctx context.Context, orgID int64, searchTerm string) ([]*gen.Client, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	orgField := query.FieldFor[gen.Client, int64](fa, "organization_id")
	nameField := query.FieldFor[gen.Client, string](fa, "name")
	emailField := query.FieldFor[gen.Client, string](fa, "email")

	orgQ := query.NewQ(orgField.Eq(orgID))
	searchQ := query.NewQ(nameField.Contains(searchTerm)).
		Or(query.NewQ(emailField.Contains(searchTerm)))

	filterQ := orgQ.And(searchQ)

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// Create creates a new client
func (r *ClientRepository) Create(ctx context.Context, client *gen.Client) error {
	return r.manager.Create(ctx, client)
}

// Update updates a client
func (r *ClientRepository) Update(ctx context.Context, client *gen.Client) error {
	return r.manager.Update(ctx, client)
}

// Delete deletes a client
func (r *ClientRepository) Delete(ctx context.Context, client *gen.Client) error {
	return r.manager.Delete(ctx, client)
}

// ============================================================================
// Skill Repository
// ============================================================================

// SkillRepository provides data access for Skill model
type SkillRepository struct {
	db      *db.DB
	manager *query.Manager[gen.Skill]
}

// NewSkillRepository creates a new skill repository
func NewSkillRepository(database *db.DB) (*SkillRepository, error) {
	manager, err := query.NewManager[gen.Skill]("skills")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &SkillRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a skill by ID
func (r *SkillRepository) GetByID(ctx context.Context, id int64) (*gen.Skill, error) {
	return r.manager.Get(ctx, id)
}

// GetByName retrieves a skill by name
func (r *SkillRepository) GetByName(ctx context.Context, name string) (*gen.Skill, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	nameField := query.FieldFor[gen.Skill, string](fa, "name")
	qs, err := r.manager.Filter(nameField.Eq(name))
	if err != nil {
		return nil, err
	}

	return qs.Get(ctx)
}

// GetAll retrieves all skills
func (r *SkillRepository) GetAll(ctx context.Context) ([]*gen.Skill, error) {
	return r.manager.All(ctx)
}

// GetActiveSkills retrieves active skills
func (r *SkillRepository) GetActiveSkills(ctx context.Context) ([]*gen.Skill, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	isActiveField := query.FieldFor[gen.Skill, bool](fa, "is_active")
	qs, err := r.manager.Filter(isActiveField.Eq(true))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("category"), query.Asc("name")).All(ctx)
}

// GetByCategory retrieves skills by category
func (r *SkillRepository) GetByCategory(ctx context.Context, category string) ([]*gen.Skill, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	categoryField := query.FieldFor[gen.Skill, string](fa, "category")
	qs, err := r.manager.Filter(categoryField.Eq(category))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("name")).All(ctx)
}

// SearchSkills searches skills by name or category
func (r *SkillRepository) SearchSkills(ctx context.Context, searchTerm string) ([]*gen.Skill, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	nameField := query.FieldFor[gen.Skill, string](fa, "name")
	categoryField := query.FieldFor[gen.Skill, string](fa, "category")

	searchQ := query.NewQ(nameField.Contains(searchTerm)).
		Or(query.NewQ(categoryField.Contains(searchTerm)))

	qs, err := r.manager.Filter(searchQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("category"), query.Asc("name")).All(ctx)
}

// Create creates a new skill
func (r *SkillRepository) Create(ctx context.Context, skill *gen.Skill) error {
	return r.manager.Create(ctx, skill)
}

// Update updates a skill
func (r *SkillRepository) Update(ctx context.Context, skill *gen.Skill) error {
	return r.manager.Update(ctx, skill)
}

// Delete deletes a skill
func (r *SkillRepository) Delete(ctx context.Context, skill *gen.Skill) error {
	return r.manager.Delete(ctx, skill)
}

// ============================================================================
// SubscriptionTier Repository
// ============================================================================

// SubscriptionTierRepository provides data access for SubscriptionTier model
type SubscriptionTierRepository struct {
	db      *db.DB
	manager *query.Manager[gen.SubscriptionTier]
}

// NewSubscriptionTierRepository creates a new subscription tier repository
func NewSubscriptionTierRepository(database *db.DB) (*SubscriptionTierRepository, error) {
	manager, err := query.NewManager[gen.SubscriptionTier]("subscription_tiers")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}
	manager.SetDB(database)

	return &SubscriptionTierRepository{
		db:      database,
		manager: manager,
	}, nil
}

// GetByID retrieves a subscription tier by ID
func (r *SubscriptionTierRepository) GetByID(ctx context.Context, id int64) (*gen.SubscriptionTier, error) {
	return r.manager.Get(ctx, id)
}

// GetBySlug retrieves a subscription tier by slug
func (r *SubscriptionTierRepository) GetBySlug(ctx context.Context, slug string) (*gen.SubscriptionTier, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	slugField := query.FieldFor[gen.SubscriptionTier, string](fa, "slug")
	qs, err := r.manager.Filter(slugField.Eq(slug))
	if err != nil {
		return nil, err
	}

	return qs.Get(ctx)
}

// GetAll retrieves all subscription tiers
func (r *SubscriptionTierRepository) GetAll(ctx context.Context) ([]*gen.SubscriptionTier, error) {
	return r.manager.All(ctx)
}

// GetActiveTiers retrieves active subscription tiers
func (r *SubscriptionTierRepository) GetActiveTiers(ctx context.Context) ([]*gen.SubscriptionTier, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	isActiveField := query.FieldFor[gen.SubscriptionTier, bool](fa, "is_active")
	qs, err := r.manager.Filter(isActiveField.Eq(true))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("monthly_price")).All(ctx)
}

// GetFeaturedTiers retrieves featured subscription tiers
func (r *SubscriptionTierRepository) GetFeaturedTiers(ctx context.Context) ([]*gen.SubscriptionTier, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	isFeaturedField := query.FieldFor[gen.SubscriptionTier, bool](fa, "is_featured")
	qs, err := r.manager.Filter(isFeaturedField.Eq(true))
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("monthly_price")).All(ctx)
}

// GetTiersByPriceRange retrieves tiers within a price range
func (r *SubscriptionTierRepository) GetTiersByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*gen.SubscriptionTier, error) {
	fa, err := r.manager.GetFieldAccessor()
	if err != nil {
		return nil, err
	}

	priceField := query.FieldFor[gen.SubscriptionTier, float64](fa, "monthly_price")

	filterQ := query.NewQ(priceField.Gte(minPrice)).
		And(query.NewQ(priceField.Lte(maxPrice)))

	qs, err := r.manager.Filter(filterQ)
	if err != nil {
		return nil, err
	}

	return qs.OrderBy(query.Asc("monthly_price")).All(ctx)
}

// Create creates a new subscription tier
func (r *SubscriptionTierRepository) Create(ctx context.Context, tier *gen.SubscriptionTier) error {
	return r.manager.Create(ctx, tier)
}

// Update updates a subscription tier
func (r *SubscriptionTierRepository) Update(ctx context.Context, tier *gen.SubscriptionTier) error {
	return r.manager.Update(ctx, tier)
}

// Delete deletes a subscription tier
func (r *SubscriptionTierRepository) Delete(ctx context.Context, tier *gen.SubscriptionTier) error {
	return r.manager.Delete(ctx, tier)
}
