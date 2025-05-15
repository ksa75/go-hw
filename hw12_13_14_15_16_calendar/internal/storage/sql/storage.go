package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // так надо
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"mycalendar/internal/storage"
)

var _ storage.BaseStorage = (*Storage)(nil) // interface assertion at compile time

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) (err error) {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

///go:embed migrations/*.sql
// var embedMigrations embed.FS

func (s *Storage) Migrate(ctx context.Context, migrate string) (err error) {
	// goose.SetBaseFS(embedMigrations)
	_ = ctx
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, e storage.Event) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO events (user_id, title, description, start_date_time, duration, notice_before, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, e.UserID, e.Title, e.Description, e.StartDateTime, e.Duration, e.NoticeBefore, e.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return storage.ErrDateBusy
		}
		return fmt.Errorf("cannot insert: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE events
		SET title = $1, description = $2, duration = $3, notice_before = $4
		WHERE user_id = $5 AND start_date_time = $6
	`, e.Title, e.Description, e.Duration, e.NoticeBefore, e.UserID, e.StartDateTime)
	if err != nil {
		return fmt.Errorf("cannot update: %w", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, userID string, start time.Time) error {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM events
		WHERE user_id = $1 AND start_date_time = $2
	`, userID, start)
	if err != nil {
		return fmt.Errorf("cannot delete: %w", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (s *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, title, description, start_date_time, duration, notice_before, created_at FROM Events
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var Events []storage.Event

	for rows.Next() {
		var e storage.Event

		if err := rows.Scan(
			&e.UserID,
			&e.Title,
			&e.Description,
			&e.StartDateTime,
			&e.Duration,
			&e.NoticeBefore,
			&e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}

		Events = append(Events, e)
	}
	return Events, rows.Err()
}

func (s *Storage) GetEventsByDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, title, description, start_date_time, duration, notice_before, created_at
		FROM Events
		WHERE start_date_time >= $1 AND start_date_time < $2
	`, date.Truncate(24*time.Hour), date.AddDate(0, 0, 1).Truncate(24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.UserID, &e.Title, &e.Description, &e.StartDateTime,
			&e.Duration, &e.NoticeBefore, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *Storage) GetEventsByWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	// Находим понедельник текущей недели
	isoWeekday := int(date.Weekday())
	if isoWeekday == 0 {
		isoWeekday = 7
	}
	startOfWeek := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, -isoWeekday+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, title, description, start_date_time, duration, notice_before, created_at
		FROM events
		WHERE start_date_time >= $1 AND start_date_time < $2
	`, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("cannot select by week: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.UserID, &e.Title, &e.Description, &e.StartDateTime,
			&e.Duration, &e.NoticeBefore, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *Storage) GetEventsByMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	year := date.Year()
	month := int(date.Month())

	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, title, description, start_date_time, duration, notice_before, created_at
		FROM Events
		WHERE EXTRACT(MONTH FROM start_date_time) = $1
		  AND EXTRACT(YEAR FROM start_date_time) = $2
	`, month, year)
	if err != nil {
		return nil, fmt.Errorf("cannot select by month: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(
			&e.UserID, &e.Title, &e.Description, &e.StartDateTime,
			&e.Duration, &e.NoticeBefore, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func isUniqueViolation(err error) bool {
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
