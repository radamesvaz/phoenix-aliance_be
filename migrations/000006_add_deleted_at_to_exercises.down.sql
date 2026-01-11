-- Remove deleted_at column from exercises table
DROP INDEX IF EXISTS idx_exercises_deleted_at;
ALTER TABLE exercises DROP COLUMN IF EXISTS deleted_at;
