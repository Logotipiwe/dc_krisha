package service

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"strconv"
	"time"
)

var Bot *tgbotapi.BotAPI

func InitBot() *tgbotapi.BotAPI {
	botToken := config.GetConfig("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go ReceiveMessages()
	return bot
}

func SendMessageInTg(msg string) {
	ownerChatIdStr := GetOwnerChatID()
	if ownerChatIdStr == "" {
		log.Println(errors.New("empty owner chat, unable to send message"))
		return
	}
	ownerChatID, _ := strconv.ParseInt(ownerChatIdStr, 10, 64)
	message := tgbotapi.NewMessage(ownerChatID, msg)
	_, err := Bot.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		acceptUserMessage(update.Message.Text)
	}
}

func acceptUserMessage(text string) {
	newInterval, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		Interval = time.Duration(newInterval) * time.Second
		SendMessageInTg("Interval set to " + Interval.String())
	} else if text == "/start" {
		Enabled = true
	} else if text == "/stop" {
		Enabled = false
	} else {
		Filters = text
	}
}
