package pkg

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	. "github.com/logotipiwe/dc_go_config_lib"
	"math"
)

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

const DateFormat = "2006-01-02 15:04:05"

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

func GetChatNameFromUpdate(update tgbotapi.Update) string {
	if update.Message.Chat.Type == "private" {
		return update.Message.Chat.UserName
	} else if update.Message.Chat.Type == "group" {
		return update.Message.Chat.Title
	} else {
		return "unknown"
	}
}
