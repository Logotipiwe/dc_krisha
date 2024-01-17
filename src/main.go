package main

import (
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"io"
	"krisha/src/http"
	service2 "krisha/src/pkg"
	"log"
	"os"
)

func init() {

}

func main() {
	setLogInFile("app.log")
	err, db := initializeApp()
	if err != nil {
		log.Fatal("Failed to initialize gorm: ", err)
	}

	services := InitServices(db)

	err = services.TgService.SendMessageToOwner("Parser started")
	if err != nil {
		panic(err)
	}

	err = services.ParserService.InitOwnerParserSettings()
	if err != nil {
		panic(err)
	}
	err = services.ParserService.StartParsersFromDb()
	if err != nil {
		panic(err)
	}
	go services.TgInteractor.StartHandleTgMessages()
	http.NewController().Start()
}

func setLogInFile(s string) {
	file, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)
}

func initializeApp() (error, *gorm.DB) {
	config.LoadDcConfigDynamically(3)
	db, err := service2.NewGormDb()
	if err != nil {
		panic(err)
	}
	return err, db
}
