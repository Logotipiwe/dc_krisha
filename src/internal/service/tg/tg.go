package tg

import (
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

func (s *TgService) StartReceiveMessages(handler func(text string) error) {
	go s.bot.ReceiveMessages(handler)
	go s.logBot.ReceiveMessages(handler)
}

func (s *TgService) SendMessageInTg(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.bot.SendMessageInTg(ownerChatID, text)
}

func (s *TgService) SendLogMessageInTg(text string) error {
	ownerChatID := pkg.GetOwnerChatID()
	return s.logBot.SendMessageInTg(ownerChatID, text)
}

func (s *TgService) SendMessageInTgWithImages(msg string, images []string) error {
	return s.bot.SendMessageInTgWithImages(pkg.GetOwnerChatID(), msg, images)
}
