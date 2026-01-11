-- Add deleted_at column to exercises table for soft delete
ALTER TABLE exercises 
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

-- Create index for filtering non-deleted exercises
CREATE INDEX IF NOT EXISTS idx_exercises_deleted_at ON exercises(deleted_at) WHERE deleted_at IS NULL;
