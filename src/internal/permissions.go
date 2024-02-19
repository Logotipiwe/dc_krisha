package internal

import (
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/pkg"
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
	return pkg.GetAutoGrantLimit() > 0
}

func (s PermissionsService) GrantAccess(chat int64) error {
	allowedChat := domain.AllowedChat{ChatID: chat}
	return s.allowedChatRepository.CreateIfNotExists(&allowedChat)
}

func (s PermissionsService) DenyAccess(chat int64) error {
	return s.allowedChatRepository.Delete(chat)
}
