package main

import (
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"io"
	"krisha/src/internal"
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

	services := internal.InitServices(db)

	err = services.TgService.SendMessageInTg("Parser started and waiting for enable...")
	if err != nil {
		panic(err)
	}
	services.TgInteractor.StartHandleTgMessages()
	services.ParserService.StartParse()
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
