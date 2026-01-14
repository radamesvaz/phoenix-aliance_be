-- Add soft delete support to workouts
ALTER TABLE workouts
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Filter indexes to exclude deleted rows
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workouts_user_id_created_at;
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_workouts_user_id_created_at ON workouts(user_id, created_at DESC) WHERE deleted_at IS NULL;

