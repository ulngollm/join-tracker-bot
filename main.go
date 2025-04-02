package main

import (
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
)

var bot *tele.Bot

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
		Poller: &tele.LongPoller{
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
	bot.Handle(tele.OnChatMember, joinHandler)

	//see also
	//tele.OnUserJoined // todo актуальность сомнительна
	//tele.OnChatJoinRequest // событие запроса на вструпление в чат. Только для тех, где установлена премодерация? see approveChatJoinRequest
	//tele.OnMyChatMember  // показывает изменение статуса бота в этом чате

	bot.Start()
}

func joinHandler(c tele.Context) error {
	if c.ChatMember().NewChatMember.Role == tele.Left {
		return nil
	}
	fmt.Printf("user %d joined to chat %s\n", c.ChatMember().Sender.ID, c.ChatMember().Chat.Title)
	return nil
}
