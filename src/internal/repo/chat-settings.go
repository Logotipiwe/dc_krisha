package repo

import (
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
)

type ParserSettingsRepository struct {
	db *gorm.DB
}

func NewParserSettingsRepository(db *gorm.DB) *ParserSettingsRepository {
	return &ParserSettingsRepository{db: db}
}

func (r *ParserSettingsRepository) Create(chat *domain.ParserSettings) error {
	return r.db.Create(chat).Error
}

func (r *ParserSettingsRepository) Get(chatID int64) (*domain.ParserSettings, error) {
	var chat domain.ParserSettings
	err := r.db.First(&chat, "chat_id = ?", chatID).Error
	return &chat, err
}

func (r *ParserSettingsRepository) Update(chat *domain.ParserSettings) error {
	return r.db.Debug().Model(&domain.ParserSettings{}).Where("chat_id = ?", chat.ID).
		Updates(chat).Error
}

func (r *ParserSettingsRepository) Delete(chatID int64) error {
	return r.db.Delete(&domain.ParserSettings{}, chatID).Error
}

func (r *ParserSettingsRepository) UpdateOrCreate(d *domain.ParserSettings) error {
	return r.db.Table(d.TableName()).Where("chat_id = ?", d.ID).Assign(d).FirstOrCreate(d).Error
}
