package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/sqlite"
	tele "gopkg.in/telebot.v3"
)

type JoinHandler struct {
	db *sql.DB
}

func NewJoinHandler(databasePath string) (*JoinHandler, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Create the join_events table if it doesn't exist
	createTableQuery := `CREATE TABLE IF NOT EXISTS join_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		chat_id INTEGER,
		chat_title TEXT
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &JoinHandler{db: db}, nil
}

func (h *JoinHandler) LogJoin(chat *tele.Chat, userID int64) error {
	insertQuery := `INSERT INTO join_events (user_id, chat_id, chat_title) VALUES (?, ?, ?)`
	_, err := h.db.Exec(insertQuery, userID, chat.ID, chat.Title)
	if err != nil {
		return fmt.Errorf("failed to log join event: %w", err)
	}

	log.Printf("Logged join event: user %d joined chat %s", userID, chat.Title)
	return nil
}
