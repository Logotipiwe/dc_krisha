package tghttp

import (
	"krisha/src/internal/service"
	"krisha/src/internal/service/tg"
	"strconv"
	"time"
)

type TgInteractor struct {
	tgService *tg.TgService
}

func NewTgInteractor(tgService *tg.TgService) *TgInteractor {
	return &TgInteractor{tgService}
}

func (i *TgInteractor) acceptUserMessage(text string) error {
	newInterval, err := strconv.ParseInt(text, 10, 64)
	if err == nil {
		service.Interval = time.Duration(newInterval) * time.Second
		return i.tgService.SendMessageInTg("Interval set to " + service.Interval.String())
	} else if text == "/start" {
		service.Enabled = true
		return i.tgService.SendMessageInTg("Parser started")
	} else if text == "/stop" {
		service.Enabled = false
		return i.tgService.SendMessageInTg("Parser stopped")
	} else {
		service.Filters = text
	}
	return nil
}

func (i *TgInteractor) StartHandleTgMessages() {
	i.tgService.StartReceiveMessages(i.acceptUserMessage)
}
