package parser

import (
	"krisha/src/internal/domain"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
)

type Factory struct {
	tgService    *tg.TgService
	krishaClient *api.KrishaClientService
}

func NewParserFactory(
	tgService *tg.TgService,
	krishaClient *api.KrishaClientService,
) *Factory {
	return &Factory{
		tgService:    tgService,
		krishaClient: krishaClient,
	}
}

func (f *Factory) CreateParser(settings *domain.ParserSettings) (*Parser, error) {
	if err := settings.ValidForStartParse(); err != nil {
		return nil, err
	}
	return newParser(settings, f), nil
}
