package web

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func init() {
	log.Println("init(): Loading telegram")
	var err error
	Bot, err = tgbotapi.NewBotAPI(telegrambotapikey)
	if err != nil {
		log.Fatalln(err)
	}
	go handleTg()
	go tgProxyMessages()
}

var messages = NewBroker()
var room = int64(-1001725432775)

func init() {
	log.Println("init(): go messages.Start()")
	go messages.Start()
}
func tgProxyMessages() {
	log.Println("tgProxyMessages(): init")
	for msg := range messages.Subscribe() {
		if strings.HasPrefix(msg, "[tg]") {
			continue
		}
		Bot.Send(tgbotapi.NewMessage(room, msg))
	}
}

func handleTg() {
	log.Println("handleTg(): Ready!")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Text == ".id" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
			msg.ReplyToMessageID = update.Message.MessageID
			Bot.Send(msg)
		}
		if update.Message.Text == ".id2" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.From.ID, 10))
			msg.ReplyToMessageID = update.Message.MessageID
			Bot.Send(msg)
		}
		if strings.HasPrefix(update.Message.Text, ".banchat") && update.Message.From.ID == 5349975097 {
			var u User
			name := update.Message.Text[len(".banchat")+1:]

			db.First(&u, "username = ?", name)
			if u.Username != name {
				Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid username"))
				return
			}
			u.ChatBanned = !u.ChatBanned
			if u.ChatBanned {
				Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "User got banned!"))
			} else {
				Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "User got unbanned!"))
			}
			db.Save(&u)
		}
		if update.Message.Chat.ID == room {
			messages.Publish("[tg]" + update.Message.From.UserName + ": " + update.Message.Text + update.Message.Caption)
		}
	}
}
