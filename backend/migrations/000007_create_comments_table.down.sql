-- Drop indexes
DROP INDEX IF EXISTS idx_comments_created;
DROP INDEX IF EXISTS idx_comments_author;
DROP INDEX IF EXISTS idx_comments_ticket;

-- Drop table
DROP TABLE IF EXISTS comments;
