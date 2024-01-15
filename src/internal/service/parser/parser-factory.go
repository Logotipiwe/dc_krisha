package parser

import (
	"krisha/src/internal/domain"
)

type Factory struct {
}

func NewParserFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateParser(settings *domain.ParserSettings) (*Parser, error) {
	if err := settings.ValidForStartParse(); err != nil {
		return nil, err
	}
	return newParser(settings), nil
}
