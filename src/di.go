//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"krisha/src/http"
	"krisha/src/internal"
	"krisha/src/internal/repo"
	"krisha/src/internal/service/admin"
	"krisha/src/internal/service/apartments"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/tghttp"
)

func InitServices(db *gorm.DB, tgServicer tg.TgServicer) *Services {
	wire.Build(
		http.NewController,
		tghttp.NewTgInteractor,
		api.NewKrishaClientService,
		apartments.NewApsTgSenderService,
		apartments.NewApsLoggerService,
		parser.NewParserService,
		repo.NewParserSettingsRepository,
		repo.NewAllowedChatRepository,
		internal.NewPermissionsService,
		parser.NewParserFactory,
		admin.NewService,
		NewServices,
	)
	return &Services{}
}
