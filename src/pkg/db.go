package pkg

import (
	"github.com/jinzhu/gorm"
)

func NewGormDb() (*gorm.DB, error) {
	err, db := InitDb()
	if err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
