DROP INDEX IF EXISTS idx_skill_category;

DROP INDEX IF EXISTS idx_skill_name;

DROP INDEX IF EXISTS idx_task_completed;

DROP INDEX IF EXISTS idx_task_status;

DROP INDEX IF EXISTS idx_task_assigned;

DROP INDEX IF EXISTS idx_task_project;

DROP INDEX IF EXISTS idx_client_email;

DROP INDEX IF EXISTS idx_client_org;

DROP INDEX IF EXISTS idx_proj_code;

DROP INDEX IF EXISTS idx_proj_status;

DROP INDEX IF EXISTS idx_proj_client;

DROP INDEX IF EXISTS idx_proj_manager;

DROP INDEX IF EXISTS idx_proj_org;

DROP INDEX IF EXISTS idx_emp_active;

DROP INDEX IF EXISTS idx_emp_number;

DROP INDEX IF EXISTS idx_emp_email;

DROP INDEX IF EXISTS idx_emp_dept;

DROP INDEX IF EXISTS idx_emp_org;

DROP INDEX IF EXISTS idx_dept_code;

DROP INDEX IF EXISTS idx_dept_manager;

DROP INDEX IF EXISTS idx_dept_parent;

DROP INDEX IF EXISTS idx_dept_org;

DROP INDEX IF EXISTS idx_tier_active;

DROP INDEX IF EXISTS idx_tier_slug;

DROP INDEX IF EXISTS idx_org_active;

DROP INDEX IF EXISTS idx_org_subscription;

DROP INDEX IF EXISTS idx_org_domain;

DROP INDEX IF EXISTS idx_org_slug;

DROP TABLE IF EXISTS skills CASCADE;

DROP TABLE IF EXISTS tasks CASCADE;

DROP TABLE IF EXISTS clients CASCADE;

DROP TABLE IF EXISTS projects CASCADE;

DROP TABLE IF EXISTS employees CASCADE;

DROP TABLE IF EXISTS departments CASCADE;

DROP TABLE IF EXISTS subscription_tiers CASCADE;

DROP TABLE IF EXISTS organizations CASCADE;

DROP TABLE IF EXISTS schema_migrations;
