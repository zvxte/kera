ALTER TABLE habit_histories
ADD CONSTRAINT habit_id_date_unique UNIQUE (habit_id, date);
