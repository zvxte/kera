INSERT INTO habit_statuses (id, name)
VALUES (0, 'Active'), (1, 'Ended')
ON CONFLICT (id) DO NOTHING;
