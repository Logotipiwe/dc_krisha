package main

import (
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/service/tg"
	service2 "krisha/src/pkg"
	"log"
)

var err, db = initializeApp()

var services = initServices(db)

func main() {
	if err != nil {
		log.Fatal("Failed to initialize gorm: ", err)
	}

	err = services.TgService.SendMessageToOwner("Parser started")
	if err != nil {
		panic(err)
	}

	//err, _ = services.ParserService.InitOwnerParserSettings()
	//if err != nil {
	//	panic(err)
	//}
	err = services.ParserService.StartParsersFromDb()
	if err != nil {
		panic(err)
	}
	go services.Controller.Start()
	services.TgInteractor.StartHandleTgMessages()
}

func initServices(db *gorm.DB) *Services {
	isTesting, _ := config.GetConfigBool("IS_TESTING")
	var tgServicer tg.TgServicer
	if isTesting {
		tgServicer = tg.NewTgMockService()
	} else {
		tgServicer = tg.NewTgService()
	}
	services := InitServices(db, tgServicer)
	return services
}

func initializeApp() (error, *gorm.DB) {
	config.LoadDcConfigDynamically(3)
	db, err := service2.NewGormDb()
	if err != nil {
		panic(err)
	}
	return err, db
}
