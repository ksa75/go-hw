-- +goose Up
ALTER TABLE events
ADD CONSTRAINT unique_user_start_time UNIQUE (user_id, start_date_time);
