CREATE TABLE IF NOT EXISTS habit_histories(
    habit_id UUID NOT NULL REFERENCES habits(id),
    date DATE NOT NULL,
    days BIGINT NOT NULL
);
