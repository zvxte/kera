CREATE TABLE IF NOT EXISTS habits(
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    status SMALLINT NOT NULL REFERENCES habit_statuses(id),
    title VARCHAR(64) NOT NULL,
    description VARCHAR(256) NOT NULL,
    tracked_week_days SMALLINT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
);
