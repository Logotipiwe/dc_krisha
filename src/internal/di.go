//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/service"
	"krisha/src/internal/service/apartments"
	"krisha/src/internal/service/api"
	"krisha/src/internal/service/tg"
	"krisha/src/tghttp"
)

func InitServices(db *gorm.DB) *Services {
	wire.Build(
		service.NewParserService,
		tg.NewTgService,
		tghttp.NewTgInteractor,
		api.NewKrishaClientService,
		apartments.NewApsTgSenderService,
		apartments.NewApsCacheService,
		apartments.NewApsLoggerService,
		NewServices,
	)
	return &Services{}
}
