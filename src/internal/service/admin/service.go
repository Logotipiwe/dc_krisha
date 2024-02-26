package admin

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/parser"
	"krisha/src/pkg"
)

type Service struct {
	parserSettingsRepo *repo.ParserSettingsRepository
	knownChatsRepo     *repo.KnownChatsRepo
}

func NewService(
	parserSettingsRepo *repo.ParserSettingsRepository,
	knownChatsRepo *repo.KnownChatsRepo,
) *Service {
	return &Service{
		parserSettingsRepo: parserSettingsRepo,
		knownChatsRepo:     knownChatsRepo,
	}
}

func (s *Service) GetGeneralInfo() (*domain.AdminInfo, error) {
	info := &domain.AdminInfo{}
	active, err := s.parserSettingsRepo.GetActive()
	if err != nil {
		return nil, err
	}
	info.ActiveParsers = active
	info.DefaultInterval = parser.DefaultIntervalSec
	info.AutoGrantLimit = pkg.GetAutoGrantLimit()
	return info, nil
}

func (s *Service) GetKnownChats() ([]*domain.KnownChat, error) {
	return s.knownChatsRepo.GetAll()
}

func (s *Service) SaveKnownChatInfo(update tgbotapi.Update) error {
	bytes, err := json.Marshal(update)
	if err != nil {
		return err
	}
	chat := &domain.KnownChat{
		ChatID:   update.Message.Chat.ID,
		ChatInfo: string(bytes),
	}
	return s.knownChatsRepo.Create(chat)
}
