package internal

import (
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
	return s.allowedChatRepository.Exists(chatID)
}

func (s PermissionsService) GrantAccess(chat int64) error {
	allowedChat := domain.AllowedChat{ChatID: chat}
	return s.allowedChatRepository.Create(&allowedChat)
}

func (s PermissionsService) DenyAccess(chat int64) error {
	return s.allowedChatRepository.Delete(chat)
}
