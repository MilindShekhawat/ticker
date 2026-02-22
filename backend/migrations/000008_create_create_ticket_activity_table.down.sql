-- Drop indexes
DROP INDEX IF EXISTS idx_activity_created;
DROP INDEX IF EXISTS idx_activity_type;
DROP INDEX IF EXISTS idx_activity_actor;
DROP INDEX IF EXISTS idx_activity_ticket;

-- Drop table
DROP TABLE IF EXISTS ticket_activity;
