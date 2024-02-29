package db_messages_log

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DbMessagesLogger interface {
	LogOutcomingMessage(chatID int64, text string)
	LogIncomingUpdate(update tgbotapi.Update)
	LogOutcomingMessageWithImages(chatID int64, text string, images []string)
}
