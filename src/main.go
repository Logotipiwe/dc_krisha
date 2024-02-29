package main

import (
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/repo"
	db_messages_log "krisha/src/internal/service/db-messages-log"
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

	err = services.ParserService.StartParsersFromDb()
	if err != nil {
		panic(err)
	}
	go services.Controller.Start()
	services.ParserService.RunLimitsChecker()
	services.TgInteractor.StartHandleTgMessages()
}

func initServices(db *gorm.DB) *Services {
	isTesting, _ := config.GetConfigBool("IS_TESTING")
	messagesLogRepo := repo.NewMessagesLogRepo(db)
	var tgServicer tg.TgServicer
	dbLogger := db_messages_log.NewLoggerService(messagesLogRepo)
	if isTesting {
		tgServicer = tg.NewTgMockService()
	} else {
		tgServicer = tg.NewTgService(dbLogger)
	}
	services := InitServices(db, tgServicer, dbLogger)
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
