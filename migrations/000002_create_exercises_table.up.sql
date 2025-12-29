-- Create exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id_exercise BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_exercises_user_id ON exercises(user_id);
CREATE INDEX IF NOT EXISTS idx_exercises_user_id_name ON exercises(user_id, name);

