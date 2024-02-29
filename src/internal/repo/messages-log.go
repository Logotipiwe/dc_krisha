package repo

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
)

type MessagesLogRepo struct {
	db *gorm.DB
}

func NewMessagesLogRepo(db *gorm.DB) *MessagesLogRepo {
	return &MessagesLogRepo{db: db}
}

func (r *MessagesLogRepo) Create(msg *domain.MessageLog) error {
	msg.ID = uuid.New().String()
	return r.db.Create(msg).Error
}
