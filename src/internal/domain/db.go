package domain

import "errors"

type ParserSettings struct {
	ID                  int64 `gorm:"column:chat_id"`
	Filters             string
	IntervalSec         int
	Limit               int `gorm:"column:aps_limit"`
	Enabled             bool
	IsGrantedExplicitly bool `gorm:"column:is_granted_explicitly"`
	ApsCount            int  `gorm:"column:curr_aps_count"`
}

func (p ParserSettings) TableName() string {
	return "parsers_settings"
}

func (p ParserSettings) ValidForStartParse() error {
	if p.ID == 0 {
		return errors.New("chat id is empty")
	}
	if p.Filters == "" || p.Filters == "?" {
		return errors.New("filters are empty")
	}
	if p.IntervalSec == 0 {
		return errors.New("interval is empty")
	}
	if !p.Enabled {
		return errors.New("enabled prop is false")
	}
	return nil
}

type KnownChat struct {
	ChatID   int64  `gorm:"column:chat_id"`
	ChatInfo string `gorm:"column:chat_info"`
}

func (c KnownChat) TableName() string {
	return "known_chats"
}

type Apartment struct {
	ID       string `gorm:"primaryKey"`
	DataJson string
}

type MessageLogDirection string

const (
	Income  MessageLogDirection = "INCOME"
	Outcome MessageLogDirection = "OUTCOME"
)

type MessageLog struct {
	ID             string
	ChatID         int64               `gorm:"column:chat_id"`
	Text           string              `gorm:"column:text"`
	Direction      MessageLogDirection `gorm:"column:direction"`
	AdditionalData string              `gorm:"column:additional_data"`
}

func (c MessageLog) TableName() string {
	return "messages_log"
}
