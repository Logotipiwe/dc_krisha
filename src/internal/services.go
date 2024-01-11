package internal

import (
	"krisha/src/internal/service"
	"krisha/src/internal/service/apartments"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
	"krisha/src/tghttp"
)

type Services struct {
	ApsCacheService     *apartments.ApsCacheService
	ApsLoggerService    *apartments.ApsLoggerService
	ApsTgSenderService  *apartments.ApsTgSenderService
	KrishaClientService *api.KrishaClientService
	ParserService       *service.ParserService
	TgService           *tg.TgService
	TgInteractor        *tghttp.TgInteractor
}

func NewServices(apsCacheService *apartments.ApsCacheService, apsLoggerService *apartments.ApsLoggerService, apsTgSenderService *apartments.ApsTgSenderService, krishaClientService *api.KrishaClientService, parserService *service.ParserService, tgService *tg.TgService, tgInteractor *tghttp.TgInteractor) *Services {
	return &Services{ApsCacheService: apsCacheService, ApsLoggerService: apsLoggerService, ApsTgSenderService: apsTgSenderService, KrishaClientService: krishaClientService, ParserService: parserService, TgService: tgService, TgInteractor: tgInteractor}
}
