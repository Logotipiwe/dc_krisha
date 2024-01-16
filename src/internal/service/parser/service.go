package parser

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/tg"
	"log"
)

type Service struct {
	parserSettingsRepo *repo.ParserSettingsRepository
	parserFactory      *Factory
	tgService          *tg.TgService
}

const (
	defaultIntervalSec = 2
)

//TODO make max allowed aps size

var parsers = make(map[int64]*Parser)

func NewParserService(
	parserSettingsRepo *repo.ParserSettingsRepository,
	tgService *tg.TgService,
	parserFactory *Factory,
) *Service {
	return &Service{
		parserSettingsRepo: parserSettingsRepo,
		tgService:          tgService,
		parserFactory:      parserFactory,
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

func (s *Service) SetFiltersAndStartParser(chatID int64, filters string) (err error, parserExisted bool) {
	settings, err := s.parserSettingsRepo.Get(chatID)
	if err != nil {
		return err, false
	}

	settings.Filters = filters
	settings.Enabled = true
	err = s.parserSettingsRepo.Update(settings)
	if err != nil {
		return err, false
	}
	parser, has := parsers[chatID]
	if has {
		parser.settings = settings
		return nil, true
	} else {
		return s.startNewParser(settings), false
	}
}

func (s *Service) startNewParser(settings *domain.ParserSettings) error {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ParserNotFoundErr
		}
		return err
	}
	parser, ok := parsers[chatID]
	if !ok {
		return domain.ParserNotFoundErr
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

func (s *Service) StartParsersFromDb() error {
	settingsFromDb, err := s.parserSettingsRepo.GetAll()
	if err != nil {
		log.Println("Failed to get parser settings from the database:", err)
		return err
	}

	for _, settings := range settingsFromDb {
		if settings.Enabled {
			parser, err := s.parserFactory.CreateParser(settings)
			if err != nil {
				s.handleParserStartErr(settings, err)
				continue
			}
			if err = parser.startParsing(); err != nil {
				s.handleParserStartErr(settings, err)
				continue
			}
			parsers[settings.ID] = parser
		}
	}
	return nil
}

func (s *Service) handleParserStartErr(settings *domain.ParserSettings, err error) {
	settingsJson, _ := json.Marshal(settings)
	s.tgService.SendLogMessageToOwner(
		"Error creating parser from db. " + string(settingsJson) + ". " + err.Error())
}
