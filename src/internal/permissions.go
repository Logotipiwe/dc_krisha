package internal

import (
	"krisha/src/internal/domain"
	"krisha/src/pkg"
)

type PermissionsService struct {
}

func NewPermissionsService() *PermissionsService {
	return &PermissionsService{}
}

func (s PermissionsService) HasAccess(settings *domain.ParserSettings) bool {
	//If auto grant - settings created firstly
	agl := pkg.GetAutoGrantLimit()
	if settings != nil {
		explicitlyZero := settings.IsGrantedExplicitly && settings.Limit == 0
		if agl > 0 {
			return !explicitlyZero
		} else {
			return settings.Limit > 0
		}
	}
	return false
}
