package repository

import (
	"errors"
	"gorm.io/gorm"
	"waroka/model"
)

type IPaymentDestinationRepository interface {
	FindById(id uint) (*model.PaymentDestination, error)
	FindByRoomId(id uint) ([]model.PaymentDestination, error)
	Create(pd *model.PaymentDestination) error
	Update(pd *model.PaymentDestination) error
}

type PaymentDestinationRepository struct {
	db *gorm.DB
}

func NewPaymentDestinationRepository(db *gorm.DB) IPaymentDestinationRepository {
	return &PaymentDestinationRepository{db: db}
}

func (r *PaymentDestinationRepository) FindById(id uint) (*model.PaymentDestination, error) {
	var pd model.PaymentDestination
	if err := r.db.Preload("PaymentDetails").First(&pd, id).Error; err != nil {
		return nil, err
	}
	return &pd, nil
}

func (r *PaymentDestinationRepository) FindByRoomId(id uint) ([]model.PaymentDestination, error) {
	var pd []model.PaymentDestination
	if err := r.db.Where("room_id = ?", id).Find(&pd).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return pd, nil
}

func (r *PaymentDestinationRepository) Create(pd *model.PaymentDestination) error {
	return r.db.Create(pd).Error
}

func (r *PaymentDestinationRepository) Update(pd *model.PaymentDestination) error {
	return r.db.Save(pd).Error
}
