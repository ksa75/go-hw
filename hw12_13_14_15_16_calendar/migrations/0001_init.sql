-- +goose Up
CREATE table events (
    id              serial primary key,
    user_id         text,
    title           text,
    description     text,
    start_date_time timestamptz,
    duration        text,
    notice_before   int,
    created_at      timestamptz not null default now()
);

-- +goose Down
drop table events;
