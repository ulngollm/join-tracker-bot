package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var (
	bot         *tele.Bot
	joinHandler *JoinHandler
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}

	botToken := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:     botToken,
		ParseMode: tele.ModeMarkdown,
		Poller:    &tele.LongPoller{
			//AllowedUpdates: []string{"chat_member"},
		},
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
}

func main() {
	joinHandler = NewJoinHandler("join_events.db")

	bot.Handle(tele.OnChatMember, handleJoin)

	//see also
	//tele.OnUserJoined // todo актуальность сомнительна
	//tele.OnChatJoinRequest // событие запроса на вструпление в чат. Только для тех, где установлена премодерация? see approveChatJoinRequest
	//tele.OnMyChatMember  // показывает изменение статуса бота в этом чате

	bot.Start()
}

func handleJoin(c tele.Context) error {
	if c.ChatMember().NewChatMember.Role == tele.Left {
		return nil
	}
	userID := c.ChatMember().Sender.ID
	if err := joinHandler.LogJoin(c.Chat(), userID); err != nil {
		log.Printf("logJoin: %s", err)
	}
	log.Printf("Logged join event: user %d joined chat %s", userID, c.Chat().Title)
	return nil
}
