-- Remove soft delete support from workouts
ALTER TABLE workouts
  DROP COLUMN IF EXISTS deleted_at;

-- Recreate non-filtered indexes
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workouts_user_id_created_at;
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_user_id_created_at ON workouts(user_id, created_at DESC);

