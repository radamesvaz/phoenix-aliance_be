-- Change workouts table: replace date column with name column
ALTER TABLE workouts 
  ADD COLUMN IF NOT EXISTS name VARCHAR(255);

-- Migrate existing data: set name based on date if name is null
UPDATE workouts 
SET name = 'Workout ' || TO_CHAR(date, 'YYYY-MM-DD')
WHERE name IS NULL;

-- Make name NOT NULL after migration
ALTER TABLE workouts 
  ALTER COLUMN name SET NOT NULL;

-- Drop the date column
ALTER TABLE workouts 
  DROP COLUMN IF EXISTS date;

-- Update indexes: remove date index, keep user_id index
DROP INDEX IF EXISTS idx_workouts_user_id_date;
CREATE INDEX IF NOT EXISTS idx_workouts_user_id_created_at ON workouts(user_id, created_at DESC);


