package main

import (
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
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

	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		log.Fatalf("DB_PATH variable is not defined")
	}
	botToken, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Fatalf("bot token is not set")
	}

	joinHandler = NewJoinHandler(dbPath)
	pref := tele.Settings{
		Token:     botToken,
		ParseMode: tele.ModeMarkdown,
		Poller: &tele.LongPoller{
			AllowedUpdates: []string{
				"chat_member",
				"message_reaction",
				"message_reaction_count",
			},
		},
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
}

func main() {
	bot.Handle(tele.OnMessageReaction, onChatReaction)
	bot.Handle(tele.OnMessageReactionCount, onChannelReactions)
	bot.Handle(tele.OnChatMember, handleJoin)

	//see also
	//tele.OnUserJoined // когда пользователь вступает в группу и фиксируется сообщение о новом пользователе
	//tele.OnChatJoinRequest // событие запроса на вструпление в чат. Только для тех, где установлена премодерация? see approveChatJoinRequest
	//tele.OnMyChatMember  // показывает изменение статуса бота в этом чате

	bot.Start()
}

func handleJoin(c tele.Context) error {
	if c.ChatMember().NewChatMember.Role == tele.Left {
		return nil
	}
	//todo понять, что пользователь реально новый
	//возможно никак. OldChatMember есть и у новых пользователей
	userID := c.ChatMember().Sender.ID
	t := c.ChatMember().Time()
	if err := joinHandler.LogJoin(c.Chat(), userID, t); err != nil {
		log.Printf("logJoin: %s", err)
	}
	log.Printf("Logged join event: user %d joined chat %s", userID, c.Chat().Title)
	return nil
}

func onChannelReactions(c tele.Context) error {
	fmt.Printf("add %d %s reactions to chat %s", c.Update().MessageReactionCount.Reactions[0].Count, c.Update().MessageReactionCount.Reactions[0].Type.Emoji, c.Update().MessageReactionCount.Chat.Title)
	return nil
}

func onChatReaction(c tele.Context) error {
	fmt.Printf("%s reaction from user %s to chat %s", c.Update().MessageReaction.NewReaction[0].Emoji, c.Update().MessageReaction.User.Username, c.Update().MessageReaction.Chat.Title)
	return nil
}
