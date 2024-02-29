package db_messages_log

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"krisha/src/internal/domain"
	"krisha/src/internal/repo"
	"strings"
)

type DbLogService struct {
	repo *repo.MessagesLogRepo
}

func NewLoggerService(
	repo *repo.MessagesLogRepo,
) *DbLogService {
	return &DbLogService{repo: repo}
}

func (d *DbLogService) LogOutcomingMessageWithImages(chatID int64, text string, images []string) {
	messageLog := &domain.MessageLog{
		ChatID:         chatID,
		Text:           text,
		Direction:      domain.Outcome,
		AdditionalData: "[IMAGES] " + strings.Join(images, ","),
	}
	d.repo.Create(messageLog)
}

func (d *DbLogService) LogOutcomingMessage(chatID int64, text string) {
	messageLog := &domain.MessageLog{
		ChatID:         chatID,
		Text:           text,
		Direction:      domain.Outcome,
		AdditionalData: "",
	}
	d.repo.Create(messageLog)
}

func (d *DbLogService) LogIncomingUpdate(update tgbotapi.Update) {
	marshal, err := json.Marshal(update)
	var addData string
	if err != nil {
		addData = err.Error()
	} else {
		addData = string(marshal)
	}
	messageLog := &domain.MessageLog{
		ChatID:         update.Message.Chat.ID,
		Text:           update.Message.Text,
		Direction:      domain.Income,
		AdditionalData: addData,
	}
	d.repo.Create(messageLog)
}
