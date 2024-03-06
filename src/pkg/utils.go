package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
	"math"
	"strconv"
)

const DateFormat = "2006-01-02 15:04:05"

func InitDb() (error, *sql.DB) {
	connectionStr := fmt.Sprintf("postgres://%v:%v@%v:5432/%v?sslmode=disable",
		GetConfig("DB_LOGIN"), GetConfig("DB_PASS"),
		GetConfig("DB_HOST"), GetConfig("DB_NAME"))
	conn, err := sql.Open("postgres", connectionStr)
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

func IsTesting() bool {
	is, _ := GetConfigBool("IS_TESTING")
	return is
}

func GetAutoGrantLimit() int {
	configInt, err := GetConfigInt("AUTO_GRANT_LIMIT")
	if err != nil {
		return 0
	}
	return configInt
}

func GetAutoStopHours() int {
	return GetConfigIntOr("AUTO_STOP_HOURS", 0)
}

func GetChatNameFromUpdate(update tgbotapi.Update) string {
	if update.Message.Chat.Type == "private" {
		return update.Message.Chat.UserName
	} else if update.Message.Chat.Type == "group" {
		return update.Message.Chat.Title
	} else {
		return "unknown"
	}
}
