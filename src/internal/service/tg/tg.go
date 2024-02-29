package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/service/db-messages-log"
	"krisha/src/pkg"
)

type TgServicer interface {
	StartReceiveMessages(handler func(update tgbotapi.Update) error)
	SendMessage(chatID int64, text string) error
	SendImgMessage(chatID int64, msg string, images []string) error
	SendMessageToOwner(text string) error
	SendLogMessageToOwner(text string) error
	SendImgMessageToOwner(msg string, images []string) error
}

type TgService struct {
	bot      *BotAPI
	logBot   *BotAPI
	dbLogger db_messages_log.DbMessagesLogger
}

func NewTgService(logger db_messages_log.DbMessagesLogger) *TgService {
	return &TgService{
		bot:      NewBotAPI(config.GetConfig("BOT_TOKEN")),
		logBot:   NewBotAPI(config.GetConfig("LOG_BOT_TOKEN")),
		dbLogger: logger,
	}
}

func (s *TgService) StartReceiveMessages(handler func(update tgbotapi.Update) error) {
	go s.bot.ReceiveMessages(handler)
	s.logBot.ReceiveMessages(handler)
}

func (s *TgService) SendMessage(chatID int64, text string) error {
	go s.dbLogger.LogOutcomingMessage(chatID, text) //TODO cover with tests
	return s.bot.SendMessageInTg(chatID, text)
}

func (s *TgService) SendImgMessage(chatID int64, msg string, images []string) error {
	go s.dbLogger.LogOutcomingMessageWithImages(chatID, msg, images) //TODO cover with tests
	return s.bot.SendMessageInTgWithImages(chatID, msg, images)
}

func (s *TgService) SendMessageToOwner(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.SendMessage(ownerChatID, text)
}

func (s *TgService) SendLogMessageToOwner(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.logBot.SendMessageInTg(ownerChatID, text)
}

func (s *TgService) SendImgMessageToOwner(msg string, images []string) error {
	return s.SendImgMessage(pkg.GetOwnerChatID(), msg, images)
}
