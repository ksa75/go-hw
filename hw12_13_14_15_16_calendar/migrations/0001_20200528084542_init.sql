-- +goose Up
CREATE table events (
    id              serial primary key,
    user_id         text,
    title           text,
    description     text,
    start_date_time timestamptz,
    duration        text,
    notice_before   int,
    created_at      timestamptz not null default now(),
);

INSERT INTO books (user_id, title, description, start_date_time, duration, notice_before)
VALUES
    ('007','ДР мамы', '', '21.12.2025', '1d','15'),
    ('006','Демо проекта', 'проект Календарь', '21.07.2025 9:00', '1h','7'),
    ('006','Установочная встреча', 'Приземление инфры из селектела', '21.12.2025 11:00', '2h','3'),

-- +goose Down
drop table events;
