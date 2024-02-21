package repo

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
	"log"
)

type AllowedChatRepository struct {
	db *gorm.DB
}

func NewAllowedChatRepository(db *gorm.DB) *AllowedChatRepository {
	return &AllowedChatRepository{db: db}
}

func (r *AllowedChatRepository) Create(chat *domain.AllowedChat) error {
	return r.db.Debug().Create(chat).Error
}

func (r *AllowedChatRepository) CreateIfNotExists(chat *domain.AllowedChat) error {
	if r.Exists(chat.ChatID) {
		log.Printf("Access for chat %v already granted", chat.ChatID)
		return nil
	}
	return r.Create(chat)
}

func (r *AllowedChatRepository) Get(chatID int) (*domain.AllowedChat, error) {
	var chat domain.AllowedChat
	err := r.db.First(&chat, chatID).Error
	return &chat, err
}

func (r *AllowedChatRepository) Delete(chatID int64) error {
	return r.db.Delete(&domain.AllowedChat{}, "chat_id = ?", chatID).Error
}

func (r *AllowedChatRepository) GetAllChatsAsArray() ([]int64, error) {
	var chats []domain.AllowedChat
	err := r.db.Find(&chats).Error
	if err != nil {
		return nil, err
	}
	chatIDs := make([]int64, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ChatID
	}
	return chatIDs, nil
}

func (r *AllowedChatRepository) Exists(chatID int64) bool {
	var chat domain.AllowedChat
	err := r.db.Debug().First(&chat, "chat_id = ?", chatID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("Error checking allowed chats table!")
		log.Println(err)
	}
	if err != nil {
		return false
	}
	return true
}
