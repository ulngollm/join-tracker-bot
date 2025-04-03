package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"welcome-bot/repository"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

var (
	bot                 *tele.Bot
	joinEventRepository repository.JoinEventRepository
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

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	joinEventRepository = repository.NewJoinEventRepository(db)

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
	bot.Handle(tele.OnMessageReactionCount, onChannelReactions) // очень медленно отслеживаются
	bot.Handle(tele.OnChatMember, handleJoin)

	bot.Handle("/stats", getAdminStatistics)

	bot.Start()
}

func handleJoin(c tele.Context) error {
	if c.ChatMember().NewChatMember.Role == tele.Left {
		return nil
	}
	//todo понять, что пользователь реально новый
	//возможно никак. OldChatMember есть и у новых пользователей
	userID := c.ChatMember().Sender.ID
	chat := c.Chat()
	//todo проверять, что пользователь впервые добавляется в чат?

	e := repository.JoinEvent{
		UserID:    userID,
		ChatID:    chat.ID,
		ChatTitle: chat.Title,
		ChatType:  string(chat.Type),
		CreatedAt: c.ChatMember().Time(),
	}
	if err := joinEventRepository.Create(e); err != nil {
		return fmt.Errorf("joinEventRepository.Create: %w", err)
	}
	log.Printf("Logged join event: user %d joined chat %s", userID, c.Chat().Title)
	return nil
}

func onChannelReactions(c tele.Context) error {
	fmt.Printf("%v\n", c.Update().MessageReactionCount)
	fmt.Printf("add %d %s reactions to chat %s\n", c.Update().MessageReactionCount.Reactions[0].Count, c.Update().MessageReactionCount.Reactions[0].Type.Emoji, c.Update().MessageReactionCount.Chat.Title)
	return nil
}

func onChatReaction(c tele.Context) error {
	r := c.Update().MessageReaction
	var t string
	if len(r.NewReaction) < len(r.OldReaction) {
		t = "remove"
	} else {
		t = "add"
	}

	reactions := map[string]bool{}
	for _, reaction := range r.OldReaction {
		reactions[reaction.Emoji] = t == "remove"
	}
	for _, reaction := range r.NewReaction {
		if t == "add" {
			if _, ok := reactions[reaction.Emoji]; !ok {
				reactions[reaction.Emoji] = true
			}
		} else {
			reactions[reaction.Emoji] = false
		}
	}
	emojis := ""
	for reaction, changed := range reactions {
		if changed {
			emojis += reaction
		}
	}
	var username string
	if r.User != nil {
		username = r.User.Username
	}
	if r.ActorChat != nil {
		username = r.ActorChat.Title
	}
	fmt.Printf("%s %s reaction from %s to chat %s\n", t, emojis, username, r.Chat.Title)
	return nil
}

// only admin can see statistics of chat
func getAdminStatistics(c tele.Context) error {
	//csv := ""
	//user_id := c.Message().Chat.Recipient()
	//select from log.db where created_at >= time.now().setTime(0,0 ) - start of day
	//iterate over row
	// select chat id and check is currenct user is admin of this chat (use c.Bot().ChatMemberOf(chat_id, user ))
	// if is admin, add this row to csv
	//if isAdmin {

	//}
	return c.Send("blabla") // send csv as file
}
