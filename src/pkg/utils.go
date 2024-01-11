package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"strconv"
)

func InitDb() (error, *sql.DB) {
	connectionStr := fmt.Sprintf("%v:%v@tcp(%v)/%v", GetConfig("DB_LOGIN"), GetConfig("DB_PASS"),
		GetConfig("DB_HOST"), GetConfig("DB_NAME"))
	conn, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return err, nil
	}
	if err := conn.Ping(); err != nil {
		println(fmt.Sprintf("Error connecting database: %s", err))
		return err, nil
	}
	println("Database connected!")
	return nil, conn
}

func GetOwnerChatID() int64 {
	ownerChatIdStr := GetConfig("OWNER_TG_CHAT_ID")
	if ownerChatIdStr == "" {
		log.Println(errors.New("empty owner chat, unable to send message"))
		return 0
	}
	ownerChatID, _ := strconv.ParseInt(ownerChatIdStr, 10, 64)
	return ownerChatID
}
