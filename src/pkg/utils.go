package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"math"
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

func Map[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, 0)
	for _, value := range input {
		result = append(result, mapper(value))
	}
	return result
}

func Min(x int, y int) int {
	return int(math.Min(float64(x), float64(y)))
}
