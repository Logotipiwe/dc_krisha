package parser

import (
	"errors"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
)

type Service struct {
	parserSettingsRepo *repo.ParserSettingsRepository
	parserFactory      *Factory
}

const (
	defaultIntervalSec = 10
)

//TODO make max allowed aps size

var parsers = make(map[int64]*Parser)

func NewParserService(
	parserSettingsRepo *repo.ParserSettingsRepository,
) *Service {
	return &Service{
		parserSettingsRepo: parserSettingsRepo,
	}
}

func (s *Service) InitParserSettings(chatID int64) error {
	parserSettings := domain.ParserSettings{
		ID:          chatID,
		IntervalSec: defaultIntervalSec,
		Enabled:     false,
		Filters:     "",
	}
	err := s.parserSettingsRepo.UpdateOrCreate(&parserSettings)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SetFiltersAndStartParser(chatID int64, filters string) error {
	settings, err := s.parserSettingsRepo.Get(chatID)
	if err != nil {
		return err
	}
	settings.Filters = filters
	settings.Enabled = true
	//settings.IntervalSec = defaultIntervalSec
	err = s.parserSettingsRepo.Update(settings)
	if err != nil {
		return err
	}
	return s.startParser(settings)
}

func (s *Service) startParser(settings *domain.ParserSettings) error {
	parser, err := s.parserFactory.CreateParser(settings)
	if err != nil {
		return err
	}
	parsers[settings.ID] = parser
	err = parser.startParsing()
	return err
}

func (s *Service) StopParser(chatID int64) error {
	settings, err := s.parserSettingsRepo.Get(chatID)
	if err != nil {
		return err
	}
	parser, ok := parsers[chatID]
	if !ok {
		return errors.New("parser not found")
	}

	settings.Enabled = false
	err = s.parserSettingsRepo.Update(settings)
	if err != nil {
		return err
	}

	parser.disable()
	delete(parsers, chatID)
	return nil
}

func (s *Service) StartParsersFromDb() {
	//TODO implement...
}
