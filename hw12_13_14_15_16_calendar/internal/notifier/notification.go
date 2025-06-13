package notifier

import "time"

type Notification struct {
	EventID int64     `json:"id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"startAt"`
	UserID  string    `json:"userId"`
}
