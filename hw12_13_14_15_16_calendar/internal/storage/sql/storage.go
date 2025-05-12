package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/ksa75/go-hw/hw12_13_14_15_16_calendar/internal/storage"
)

var _ storage.BaseStorage = (*Storage)(nil) // interface assertion at compile time

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	s.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

///go:embed migrations/*.sql
//var embedMigrations embed.FS

func (s *Storage) Migrate(ctx context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (r *Storage) GetEvents(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, start_date_time, duration, description, user_id, notice_before, created_at FROM Events
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var Events []storage.Event

	for rows.Next() {
		var e storage.Event

		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.Title,
			&b.Description,
			&b.StartDateTime,
			&b.Duration,
			&b.UserID,
			&b.NoticeBefore,
			&b.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}

		Events = append(Events, e)
	}
	return Events, rows.Err()
}
