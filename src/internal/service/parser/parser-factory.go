package parser

import (
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
)

type Factory struct {
	tgService          tg.TgServicer
	krishaClient       *api.KrishaClientService
	parserSettingsRepo *repo.ParserSettingsRepository
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
	return newParser(settings, apsInFilter, f), nil
}
