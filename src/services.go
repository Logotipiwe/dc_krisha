package main

import (
	"krisha/src/internal"
	"krisha/src/internal/service/apartments"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/tghttp"
)

type Services struct {
	ApsLoggerService    *apartments.ApsLoggerService
	ApsTgSenderService  *apartments.ApsTgSenderService
	KrishaClientService *api.KrishaClientService
	TgService           *tg.TgService
	TgInteractor        *tghttp.TgInteractor
	ParserService       *parser.Service
	PermissionsService  *internal.PermissionsService
	ParserFactory       *parser.Factory
}

func NewServices(
	apsLoggerService *apartments.ApsLoggerService,
	apsTgSenderService *apartments.ApsTgSenderService,
	krishaClientService *api.KrishaClientService,
	parserSerivce *parser.Service,
	tgService *tg.TgService,
	tgInteractor *tghttp.TgInteractor,
	parserFactory *parser.Factory,
) *Services {
	return &Services{
		ApsLoggerService:    apsLoggerService,
		ApsTgSenderService:  apsTgSenderService,
		KrishaClientService: krishaClientService,
		TgService:           tgService,
		TgInteractor:        tgInteractor,
		ParserFactory:       parserFactory,
		ParserService:       parserSerivce,
	}
}
