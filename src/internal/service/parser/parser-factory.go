package parser

import (
	"krisha/src/internal/domain"
	"krisha/src/internal/service/tg"
)

type Factory struct {
	tgService *tg.TgService
}

func NewParserFactory(
	tgService *tg.TgService,
) *Factory {
	return &Factory{
		tgService: tgService,
	}
}

func (f *Factory) CreateParser(settings *domain.ParserSettings) (*Parser, error) {
	if err := settings.ValidForStartParse(); err != nil {
		return nil, err
	}
	return newParser(settings, f), nil
}
