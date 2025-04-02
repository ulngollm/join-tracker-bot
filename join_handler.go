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

func NewJoinHandler(databasePath string) *JoinHandler {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS join_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		chat_id INTEGER,
		chat_title TEXT
	)`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("failed to create table: %s", err.Error())
	}

	return &JoinHandler{db: db}
}

func (h *JoinHandler) LogJoin(chat *tele.Chat, userID int64) error {
	q := `INSERT INTO join_events (user_id, chat_id, chat_title) VALUES (?, ?, ?)`
	_, err := h.db.Exec(q, userID, chat.ID, chat.Title)
	if err != nil {
		return fmt.Errorf("failed to log join event: %w", err)
	}

	return nil
}
