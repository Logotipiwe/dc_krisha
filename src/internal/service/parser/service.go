package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
	"krisha/src/pkg"
	"log"
	"strconv"
)

type Service struct {
	ParserSettingsRepo *repo.ParserSettingsRepository
	parserFactory      *Factory
	tgService          tg.TgServicer
	krishaClient       *api.KrishaClientService
}

var DefaultIntervalSec = 120

var parsers = make(map[int64]*Parser)

func NewParserService(
	parserSettingsRepo *repo.ParserSettingsRepository,
	tgService tg.TgServicer,
	parserFactory *Factory,
	krishaClient *api.KrishaClientService,
) *Service {
	defaultIntervalSec, intErr := config.GetConfigInt("DEFAULT_INTERVAL")
	if intErr == nil && defaultIntervalSec > 0 {
		DefaultIntervalSec = defaultIntervalSec
	}
	return &Service{
		ParserSettingsRepo: parserSettingsRepo,
		tgService:          tgService,
		parserFactory:      parserFactory,
		krishaClient:       krishaClient,
	}
}

func (s *Service) InitOwnerParserSettings() (error, *domain.ParserSettings) {
	parserSettings := domain.ParserSettings{
		ID:                  pkg.GetOwnerChatID(),
		IntervalSec:         DefaultIntervalSec,
		Enabled:             false,
		Limit:               20000,
		Filters:             "",
		IsGrantedExplicitly: true,
	}
	return s.ParserSettingsRepo.UpdateOrCreate(&parserSettings), &parserSettings
}

func (s *Service) CreateParserSettingsFromExplicitGrant(chatID int64, limit int) error {
	parserSettings := domain.ParserSettings{
		ID:                  chatID,
		IntervalSec:         DefaultIntervalSec,
		Enabled:             false,
		Limit:               limit,
		Filters:             "",
		IsGrantedExplicitly: true,
	}
	return s.ParserSettingsRepo.UpdateOrCreate(&parserSettings)
}

func (s *Service) CreateParserSettingsFromAutoGrant(chatID int64) error {
	parserSettings := domain.ParserSettings{
		ID:                  chatID,
		IntervalSec:         DefaultIntervalSec,
		Enabled:             false,
		Limit:               0,
		Filters:             "",
		IsGrantedExplicitly: false,
	}
	return s.ParserSettingsRepo.UpdateOrCreate(&parserSettings)
}

func (s *Service) UpdateLimitExplicitly(settings *domain.ParserSettings, limit int) (err error, stopped bool) {
	settings.Limit = limit
	settings.IsGrantedExplicitly = true
	err = s.ParserSettingsRepo.Update(settings)
	if err != nil {
		return err, false
	}
	err, _ = s.checkLimits(settings)
	if err != nil {
		_, has := parsers[settings.ID]
		if has {
			return s.StopParser(settings.ID), true
		}
	}
	s.tgService.SendMessage(settings.ID, fmt.Sprintf("Ваш лимит изменен на %v квартир", limit))
	return nil, false
}

func (s *Service) SetFilters(chatID int64, filters string) (*domain.ParserSettings, error) {
	if chatID == pkg.GetOwnerChatID() {
		err, _ := s.InitOwnerParserSettings()
		if err != nil {
			wrapper := errors.New("error creating settings for admin")
			wrapper = errors.Join(err, wrapper)
			return nil, wrapper
		}
	}
	settings, err := s.ParserSettingsRepo.Get(chatID)
	if err != nil {
		return nil, err
	}
	settings.Filters = filters
	err = s.ParserSettingsRepo.Update(settings)
	return settings, err
}

func (s *Service) StartParser(settings *domain.ParserSettings, restartIfExists bool) (err error, existed bool) {
	_, has := parsers[settings.ID]
	if has {
		if !restartIfExists {
			return nil, true
		} else {
			err := s.StopParser(settings.ID)
			if err != nil {
				return err, true
			}
		}
	}
	err, apsCount := s.checkLimits(settings)
	if err != nil {
		return err, false
	}
	settings.Enabled = true
	err = s.ParserSettingsRepo.Update(settings)
	if err != nil {
		return err, false
	}
	return s.startNewParser(settings, apsCount), false
}

func (s *Service) checkLimits(settings *domain.ParserSettings) (err error, apsCount int) {
	mapData := s.krishaClient.RequestMapData(settings.Filters)
	apsCount = mapData.NbTotal
	var limit int
	if pkg.GetAutoGrantLimit() > 0 {
		limit = pkg.GetAutoGrantLimit()
	} else {
		limit = settings.Limit
	}
	if apsCount > limit {
		return domain.LimitExceededErr, apsCount
	}
	return nil, apsCount
}

func (s *Service) startNewParser(settings *domain.ParserSettings, apsInFilter int) error {
	parser, err := s.parserFactory.CreateParser(settings, apsInFilter)
	if err != nil {
		return err
	}
	parsers[settings.ID] = parser
	err = parser.startParsing()
	return err
}

func (s *Service) StopParser(chatID int64) error {
	settings, err := s.ParserSettingsRepo.Get(chatID)
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
	err = s.ParserSettingsRepo.Update(settings)
	if err != nil {
		return err
	}

	parser.disable()
	delete(parsers, chatID)
	s.tgService.SendLogMessageToOwner("Parser stopped for chat " + strconv.FormatInt(chatID, 10))
	return nil
}

func (s *Service) StartParsersFromDb() error {
	settingsFromDb, err := s.ParserSettingsRepo.GetAll()
	if err != nil {
		log.Println("Failed to get parser settings from the database:", err)
		return err
	}

	for _, settings := range settingsFromDb {
		if settings.Enabled {
			err, _ := s.StartParser(settings, false)
			if err != nil {
				s.handleParserStartErr(settings, err)
			}
		}
	}
	return nil
}

func (s *Service) handleParserStartErr(settings *domain.ParserSettings, err error) {
	settingsJson, _ := json.Marshal(settings)
	s.tgService.SendLogMessageToOwner(
		"Error creating parser from db. " + string(settingsJson) + ". " + err.Error())
}

func (s *Service) GetSettings(chatID int64) (settings *domain.ParserSettings, isErrNotFound bool, err error) {
	settings, err = s.ParserSettingsRepo.Get(chatID)
	if errors.Is(err, gorm.ErrRecordNotFound) && chatID == pkg.GetOwnerChatID() {
		err, parserSettings := s.InitOwnerParserSettings()
		return parserSettings, false, err
	}
	return settings, errors.Is(err, gorm.ErrRecordNotFound), err
}

func (s *Service) StopAllParsersOnlyInGoroutines() error {
	availableParsers := make(map[int64]*Parser)
	for chatID, parser := range parsers {
		availableParsers[chatID] = parser
	}
	for chatID, _ := range availableParsers {
		err := s.StopParser(chatID)
		if err != nil {
			return err
		}
	}
	return nil
}
