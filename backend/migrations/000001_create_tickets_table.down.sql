-- Drop indexes first
DROP INDEX IF EXISTS idx_tickets_project_status;
DROP INDEX IF EXISTS idx_tickets_deleted_at;
DROP INDEX IF EXISTS idx_tickets_created_at;
DROP INDEX IF EXISTS idx_tickets_creator;
DROP INDEX IF EXISTS idx_tickets_assignee;
DROP INDEX IF EXISTS idx_tickets_status;
DROP INDEX IF EXISTS idx_tickets_project;

-- Drop table
DROP TABLE IF EXISTS tickets;
