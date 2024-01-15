package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/pkg"
)

type TgService struct {
	bot    *BotAPI
	logBot *BotAPI
}

func NewTgService() *TgService {
	return &TgService{
		bot:    NewBotAPI(config.GetConfig("BOT_TOKEN")),
		logBot: NewBotAPI(config.GetConfig("LOG_BOT_TOKEN")),
	}
}

func (s *TgService) StartReceiveMessages(handler func(update tgbotapi.Update) error) {
	go s.bot.ReceiveMessages(handler)
	s.logBot.ReceiveMessages(handler)
}

func (s *TgService) SendMessage(chatID int64, text string) error {
	return s.bot.SendMessageInTg(chatID, text)
}

func (s *TgService) SendLogMessage(chatID int64, text string) error {
	return s.logBot.SendMessageInTg(chatID, text)
}

func (s *TgService) SendImgMessage(chatID int64, msg string, images []string) error {
	return s.bot.SendMessageInTgWithImages(chatID, msg, images)
}

func (s *TgService) SendMessageToOwner(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.bot.SendMessageInTg(ownerChatID, text)
}

func (s *TgService) SendLogMessageToOwner(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.logBot.SendMessageInTg(ownerChatID, text)
}

func (s *TgService) SendImgMessageToOwner(msg string, images []string) error {
	return s.bot.SendMessageInTgWithImages(pkg.GetOwnerChatID(), msg, images)
}
