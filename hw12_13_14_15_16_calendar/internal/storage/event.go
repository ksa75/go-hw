package storage

type EventsStorage interface {
	GetEvents(ctx context.Context) ([]Event, error)
}

type BaseStorage interface {
	Connect(ctx context.Context, dsn string) error
	Migrate(ctx context.Context, migrate string) error
	Close() error
	EventsStorage
}


type Event struct {
	ID    string
	UserID    string
	Title string
	Description string
	StartDateTime time.Time
	Duration string
	NoticeBefore  string
}
