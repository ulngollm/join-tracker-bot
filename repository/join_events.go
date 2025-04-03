package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/glebarez/sqlite"
)

type JoinEvent struct {
	ID        int64
	UserID    int64
	ChatID    int64
	ChatTitle string
	ChatType  string
	CreatedAt time.Time
}

type JoinEventRepository struct {
	db *sql.DB
}

func NewJoinEventRepository(db *sql.DB) JoinEventRepository {
	createTableQuery := `CREATE TABLE IF NOT EXISTS join_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		chat_id INTEGER,
		chat_title TEXT,
		chat_type TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("failed to create table: %s", err)
	}
	return JoinEventRepository{db: db}
}

func (r *JoinEventRepository) Create(event JoinEvent) error {
	_, err := r.db.Exec("INSERT INTO join_events (user_id, chat_id, chat_title, chat_type, created_at) VALUES (?, ?, ?, ?, ?)", event.UserID, event.ChatID, event.ChatTitle, event.ChatType, event.CreatedAt)
	return err
}

func (r *JoinEventRepository) GetAll() ([]*JoinEvent, error) {
	rows, err := r.db.Query("SELECT id, user_id, chat_id, chat_title, chat_type, created_at FROM join_events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*JoinEvent
	for rows.Next() {
		var event JoinEvent
		err = rows.Scan(&event.ID, &event.UserID, &event.ChatID, &event.ChatTitle, &event.ChatType, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (r *JoinEventRepository) GetFromDate(date time.Time) ([]JoinEvent, error) {
	rows, err := r.db.Query("SELECT id, user_id, chat_id, chat_title, chat_type, created_at FROM join_events WHERE created_at >= ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []JoinEvent
	for rows.Next() {
		var event JoinEvent
		err = rows.Scan(&event.ID, &event.UserID, &event.ChatID, &event.ChatTitle, &event.ChatType, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
