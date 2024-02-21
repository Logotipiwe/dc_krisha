package admin

import (
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/parser"
	"krisha/src/pkg"
)

type Service struct {
	parserSettingsRepo *repo.ParserSettingsRepository
}

func NewService(parserSettingsRepo *repo.ParserSettingsRepository) *Service {
	return &Service{parserSettingsRepo: parserSettingsRepo}
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
