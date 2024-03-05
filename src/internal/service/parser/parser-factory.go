package parser

import (
	"errors"
	"github.com/Logotipiwe/krisha_model/model"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
	"krisha/src/pkg"
	"time"
)

type Factory struct {
	tgService          tg.TgServicer
	krishaClient       *api.KrishaClientService
	parserSettingsRepo *repo.ParserSettingsRepository
	parserService      *Service
}

func NewParserFactory(
	tgService tg.TgServicer,
	krishaClient *api.KrishaClientService,
	parserSettingsRepo *repo.ParserSettingsRepository,
) *Factory {
	return &Factory{
		tgService:          tgService,
		krishaClient:       krishaClient,
		parserSettingsRepo: parserSettingsRepo,
	}
}

func (f *Factory) CreateParser(settings *domain.ParserSettings, apsInFilter int) (*Parser, error) {
	if err := settings.ValidForStartParse(); err != nil {
		return nil, err
	}

	startTime := getTimeFromSettings(settings)
	if startTime == nil {
		return nil, errors.New("wrong time passed to starting parser. time: " + settings.StartTime)
	}

	return newParser(settings, apsInFilter, *startTime, f), nil
}

func newParser(settings *domain.ParserSettings, apsInFilter int, startTime time.Time, factory *Factory) *Parser {
	return &Parser{
		factory:                    factory,
		settings:                   settings,
		areAllCurrentApsCollected:  false,
		areCollectApsTriesExceeded: false,
		enabled:                    true,
		collectedAps:               make(map[string]*model.Ap),
		initialApsCountInFilter:    apsInFilter,
		startTime:                  startTime,
	}
}

func getTimeFromSettings(settings *domain.ParserSettings) *time.Time {
	if settings.StartTime == "" {
		return nil
	}
	parse, err := time.Parse(pkg.DateFormat, settings.StartTime)
	if err != nil {
		return nil
	}
	return &parse
}
