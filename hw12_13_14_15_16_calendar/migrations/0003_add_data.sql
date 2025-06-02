-- +goose Up
INSERT INTO events (user_id, title, description, start_date_time, duration, notice_before)
VALUES
    ('007','test event1', 'test event', '2025-12-21 00:00:00', '1d', '15'),
    ('006','test event2', 'test event', '2025-07-21 09:00:00', '1h', '7'),
    ('006','test event3', 'test event', '2025-12-21 11:00:00', '2h', '3');

