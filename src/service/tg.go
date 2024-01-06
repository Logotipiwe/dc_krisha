package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"strconv"
	"time"
)

var Bot *tgbotapi.BotAPI
var LogBot *tgbotapi.BotAPI

func InitBot() *tgbotapi.BotAPI {
	botToken := config.GetConfig("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go ReceiveMessages(bot)
	return bot
}

func InitLogBot() *tgbotapi.BotAPI {
	botToken := config.GetConfig("LOG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go ReceiveMessages(bot)
	return bot
}

func SendMessageInTgWithImages(msg string, images []string) {
	ownerChatID := GetOwnerChatID()
	var message tgbotapi.Chattable
	message = tgbotapi.NewMessage(ownerChatID, msg)
	var photos = make([]interface{}, 0)
	for _, url := range images {
		photos = append(photos, tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(url)))
	}
	message = tgbotapi.NewMediaGroup(ownerChatID, photos)
	_, err := Bot.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
	SendMessageInTg(msg)
}

func SendMessageInTg(msg string) {
	ownerChatID := GetOwnerChatID()

	message := tgbotapi.NewMessage(ownerChatID, msg)
	_, err := Bot.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
}

func SendLogInTg(msg string) {
	ownerChatID := GetOwnerChatID()
	message := tgbotapi.NewMessage(ownerChatID, msg)
	_, err := LogBot.Send(message)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveMessages(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != GetOwnerChatID() {
			log.Printf("Received message from unauthorized chat: %d", update.Message.Chat.ID)
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
		SendMessageInTg("Parser started")
	} else if text == "/stop" {
		Enabled = false
		SendMessageInTg("Parser stopped")
	} else {
		Filters = text
	}
}
