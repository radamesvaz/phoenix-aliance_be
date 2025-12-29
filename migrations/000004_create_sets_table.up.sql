-- Create sets table
CREATE TABLE IF NOT EXISTS sets (
    id_set BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workouts(id_workout) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id_exercise) ON DELETE CASCADE,
    weight DECIMAL(10, 2) NOT NULL CHECK (weight >= 0),
    reps INTEGER NOT NULL CHECK (reps > 0),
    rest_seconds INTEGER CHECK (rest_seconds >= 0),
    notes TEXT,
    rpe INTEGER CHECK (rpe >= 1 AND rpe <= 10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_sets_workout_id ON sets(workout_id);
CREATE INDEX IF NOT EXISTS idx_sets_exercise_id ON sets(exercise_id);
CREATE INDEX IF NOT EXISTS idx_sets_exercise_id_created_at ON sets(exercise_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_sets_workout_exercise ON sets(workout_id, exercise_id);
