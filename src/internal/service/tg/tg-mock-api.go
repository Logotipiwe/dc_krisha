package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"krisha/src/pkg"
)

type SentMockMessage struct {
	ChatID int64
	Text   string
	Images []string
}

var SentMockMessages []SentMockMessage
var SentMockLogMessages []SentMockMessage

func GetSentMessages() []SentMockMessage {
	ans := SentMockMessages
	SentMockMessages = []SentMockMessage{}
	return ans
}

func GetSentLogMessages() []SentMockMessage {
	ans := SentMockLogMessages
	SentMockLogMessages = []SentMockMessage{}
	return ans
}

type TgMockService struct {
}

func NewTgMockService() *TgMockService {
	return &TgMockService{}
}

func (t TgMockService) StartReceiveMessages(handler func(update tgbotapi.Update) error) {
	for true {

	}
}

func (t TgMockService) SendMessage(chatID int64, text string) error {
	fmt.Printf("[MOCK] Sent message to id %v: %v\n", chatID, text)
	SentMockMessages = append(SentMockMessages, SentMockMessage{
		ChatID: chatID,
		Text:   text,
	})
	return nil
}

func (t TgMockService) SendImgMessage(chatID int64, msg string, images []string) error {
	fmt.Printf("[MOCK] Sent message to id %v: %v; with imgs: %v\n", chatID, msg, images)
	SentMockMessages = append(SentMockMessages, SentMockMessage{
		ChatID: chatID,
		Text:   msg,
		Images: images,
	})
	return nil
}

func (t TgMockService) SendMessageToOwner(text string) error {
	fmt.Printf("[MOCK] Sent message to owner: %v\n", text)
	SentMockMessages = append(SentMockMessages, SentMockMessage{
		ChatID: pkg.GetOwnerChatID(),
		Text:   text,
	})
	return nil
}

func (t TgMockService) SendLogMessageToOwner(text string) error {
	fmt.Printf("[MOCK] Sent log message to owner: %v\n", text)
	SentMockLogMessages = append(SentMockLogMessages, SentMockMessage{
		ChatID: pkg.GetOwnerChatID(),
		Text:   text,
	})
	return nil
}

func (t TgMockService) SendImgMessageToOwner(msg string, images []string) error {
	fmt.Printf("[MOCK] Sent log message to owner: %v; with images: %v\n", msg, images)
	SentMockMessages = append(SentMockMessages, SentMockMessage{
		ChatID: pkg.GetOwnerChatID(),
		Text:   msg,
		Images: images,
	})
	return nil
}
