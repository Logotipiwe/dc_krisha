package repo

import (
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
)

type KnownChatsRepo struct {
	db *gorm.DB
}

func NewKnownChatsRepo(db *gorm.DB) *KnownChatsRepo {
	return &KnownChatsRepo{db: db}
}

func (r *KnownChatsRepo) Get(chatID int64) (*domain.KnownChat, error) {
	var chat domain.KnownChat
	err := r.db.First(&chat, "chat_id = ?", chatID).Error
	return &chat, err
}

func (r *KnownChatsRepo) Update(chat *domain.KnownChat) error {
	//TODO get rid of columns listing
	return r.db.Debug().Model(chat).
		Where("chat_id = ?", chat.ChatID).
		UpdateColumns(map[string]interface{}{
			"chat_id":   chat.ChatID,
			"chat_info": chat.ChatInfo,
		},
		).Error
}

func (r *KnownChatsRepo) Delete(chatID int64) error {
	return r.db.Delete(&domain.KnownChat{}, chatID).Error
}

func (r *KnownChatsRepo) Create(d *domain.KnownChat) error {
	return r.db.Table(d.TableName()).Where("chat_id = ?", d.ChatID).Assign(d).FirstOrCreate(d).Error
}

func (r *KnownChatsRepo) GetAll() ([]*domain.KnownChat, error) {
	var chats []*domain.KnownChat
	err := r.db.Find(&chats).Error
	return chats, err
}
