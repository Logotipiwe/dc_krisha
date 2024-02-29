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
	db_messages_log "krisha/src/internal/service/db-messages-log"
	"krisha/src/internal/service/parser"
	"krisha/src/internal/service/tg"
	"krisha/src/tghttp"
)

func InitServices(
	db *gorm.DB,
	tgServicer tg.TgServicer,
	logger db_messages_log.DbMessagesLogger,
) *Services {
	wire.Build(
		http.NewController,
		tghttp.NewTgInteractor,
		api.NewKrishaClientService,
		apartments.NewApsTgSenderService,
		apartments.NewApsLoggerService,
		parser.NewParserService,
		repo.NewParserSettingsRepository,
		repo.NewKnownChatsRepo,
		internal.NewPermissionsService,
		parser.NewParserFactory,
		admin.NewService,
		NewServices,
	)
	return &Services{}
}
