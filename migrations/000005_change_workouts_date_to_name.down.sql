-- Revert: change workouts table back to date column
ALTER TABLE workouts 
  ADD COLUMN IF NOT EXISTS date TIMESTAMP WITH TIME ZONE;

-- Migrate existing data: set date based on created_at if date is null
UPDATE workouts 
SET date = created_at
WHERE date IS NULL;

-- Make date NOT NULL after migration
ALTER TABLE workouts 
  ALTER COLUMN date SET NOT NULL;

-- Drop the name column
ALTER TABLE workouts 
  DROP COLUMN IF EXISTS name;

-- Restore original indexes
DROP INDEX IF EXISTS idx_workouts_user_id_created_at;
CREATE INDEX IF NOT EXISTS idx_workouts_user_id_date ON workouts(user_id, date DESC);


