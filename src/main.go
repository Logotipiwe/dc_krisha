package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	config "github.com/logotipiwe/dc_go_config_lib"
	"io"
	"krisha/src/service"
	"log"
	"os"
)

var db *sql.DB

func init() {
	err, dbNew := initializeApp()
	if err != nil {
		panic(err)
	}
	db = dbNew
}

func main() {
	setLogInFile("app.log")
	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		log.Fatal("Failed to initialize gorm: ", err)
	}

	service.SendMessageInTg("Parser started and waiting for enable...")
	service.StartParse(gormDB)
}

func setLogInFile(s string) {
	file, err := os.OpenFile(s, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multi)
}

func initializeApp() (error, *sql.DB) {
	config.LoadDcConfigDynamically(3)
	err, db := service.InitDb()
	if err != nil {
		panic(err)
	}
	service.Bot = service.InitBot()
	service.LogBot = service.InitLogBot()
	return err, db
}
