package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	. "krisha/src/pkg"
	"log"
)

type BotAPI struct {
	botAPI *tgbotapi.BotAPI
}

func NewBotAPI(botToken string) *BotAPI {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &BotAPI{botAPI: bot}
}

func (b *BotAPI) SendMessageInTgWithImages(chatID int64, msg string, images []string) error {
	var message tgbotapi.Chattable
	message = tgbotapi.NewMessage(chatID, msg)
	var photos = make([]interface{}, 0)
	for _, url := range images {
		photos = append(photos, tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(url)))
	}
	message = tgbotapi.NewMediaGroup(chatID, photos)
	_, err := b.botAPI.Send(message)
	//ignore lib error
	if err != nil && err.Error() != "json: cannot unmarshal array into Go value of type tgbotapi.Message" {
		log.Println(err)
		return err
	}
	return b.SendMessageInTg(chatID, msg)
}

func (b *BotAPI) SendMessageInTg(chatID int64, msg string) error {
	message := tgbotapi.NewMessage(chatID, msg)
	_, err := b.botAPI.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (b *BotAPI) ReceiveMessages(handler func(update tgbotapi.Update) error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.botAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != GetOwnerChatID() {
			log.Printf("Received message from unauthorized chat: %d", update.Message.Chat.ID)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		err := handler(update)
		if err != nil {
			log.Println(err)
		}
	}
}
