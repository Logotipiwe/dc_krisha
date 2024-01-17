package internal

import (
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
)

type PermissionsService struct {
	allowedChatRepository *repo.AllowedChatRepository
}

func NewPermissionsService(
	allowedChatRepository *repo.AllowedChatRepository,
) *PermissionsService {
	return &PermissionsService{
		allowedChatRepository: allowedChatRepository,
	}
}

func (s PermissionsService) HasAccess(chatID int64) bool {
	allowed := s.allowedChatRepository.Exists(chatID)
	if allowed {
		return true
	}
	configInt, intErr := config.GetConfigInt("AUTO_GRANT_LIMIT")
	if intErr != nil {
		return false
	}
	return configInt > 0
}

func (s PermissionsService) GrantAccess(chat int64) error {
	allowedChat := domain.AllowedChat{ChatID: chat}
	return s.allowedChatRepository.CreateIfNotExists(&allowedChat)
}

func (s PermissionsService) DenyAccess(chat int64) error {
	return s.allowedChatRepository.Delete(chat)
}
