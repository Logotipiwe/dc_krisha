package repo

import (
	"github.com/jinzhu/gorm"
	"krisha/src/internal/domain"
)

type ApartmentRepository struct {
	db *gorm.DB
}

func NewApartmentRepository(db *gorm.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

func (r *ApartmentRepository) Create(apartment *domain.Apartment) error {
	return r.db.Create(apartment).Error
}

func (r *ApartmentRepository) Get(apartmentID string) (*domain.Apartment, error) {
	var apartment domain.Apartment
	err := r.db.First(&apartment, apartmentID).Error
	return &apartment, err
}

func (r *ApartmentRepository) Update(apartment *domain.Apartment) error {
	return r.db.Model(apartment).Updates(apartment).Error
}

func (r *ApartmentRepository) Delete(apartmentID string) error {
	return r.db.Delete(&domain.Apartment{}, apartmentID).Error
}
