package internal

import (
	"fmt"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/pkg"
)

type PermissionsService struct {
	allowedChatRepository *repo.AllowedChatRepository
	parserSettingsRepo    *repo.ParserSettingsRepository
}

func NewPermissionsService(
	allowedChatRepository *repo.AllowedChatRepository,
	parserSettingsRepo *repo.ParserSettingsRepository,
) *PermissionsService {
	return &PermissionsService{
		allowedChatRepository: allowedChatRepository,
		parserSettingsRepo:    parserSettingsRepo,
	}
}

func (s PermissionsService) HasAccess(settings *domain.ParserSettings) bool {
	allowed := s.allowedChatRepository.Exists(settings.ID)
	if allowed {
		return true
	}
	if pkg.GetAutoGrantLimit() > 0 {
		if settings.IsGrantedExplicitly && settings.Limit == 0 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func (s PermissionsService) GrantAccess(chat int64) error {
	allowedChat := domain.AllowedChat{ChatID: chat}
	return s.allowedChatRepository.CreateIfNotExists(&allowedChat)
}

func (s PermissionsService) DenyAccess(settings *domain.ParserSettings) error {
	err := s.allowedChatRepository.Delete(settings.ID)
	if err != nil {
		return err
	}
	settings.IsGrantedExplicitly = true
	settings.Limit = 0
	err = s.parserSettingsRepo.Update(settings)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
