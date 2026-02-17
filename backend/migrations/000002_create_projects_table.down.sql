-- Drop indexes first
DROP INDEX IF EXISTS idx_projects_deleted_at;
DROP INDEX IF EXISTS idx_projects_created_at;
DROP INDEX IF EXISTS idx_projects_creator;

-- Drop table
DROP TABLE IF EXISTS projects;
