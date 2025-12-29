-- Create workouts table
CREATE TABLE IF NOT EXISTS workouts (
    id_workout BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);
CREATE INDEX IF NOT EXISTS idx_workouts_user_id_date ON workouts(user_id, date DESC);

