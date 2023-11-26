package repository

import (
	"gorm.io/gorm"
	"waroka/model"
)

type IPaymentDetailRepository interface {
	FindById(paymentDestinationId uint) (*model.PaymentDetail, error)
	FindByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error)
	FindInCompletedByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error)
	FindActiveByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error)
	FindByUserId(id uint) ([]model.PaymentDetail, error)
	Create(pd *model.PaymentDetail) error
	Update(pd *model.PaymentDetail) error
	Delete(pd *model.PaymentDetail) error
}

type PaymentDetailRepository struct {
	db *gorm.DB
}

func NewPaymentDetailRepository(db *gorm.DB) IPaymentDetailRepository {
	return &PaymentDetailRepository{db: db}
}

func (r PaymentDetailRepository) FindById(id uint) (*model.PaymentDetail, error) {
	var pd model.PaymentDetail
	// by id のみ
	if err := r.db.First(&pd, id).Error; err != nil {
		return nil, err
	}
	return &pd, nil
}

func (r PaymentDetailRepository) FindByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error) {
	var pds []model.PaymentDetail
	// by payment_destination_id
	if err := r.db.Where("payment_destination_id = ?", paymentDestinationId).Find(&pds).Error; err != nil {
		return nil, err
	}
	return pds, nil
}

// FindActiveByPaymentDestinationId まだコンプリートしていない支払いを取得する
func (r PaymentDetailRepository) FindInCompletedByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error) {
	var pds []model.PaymentDetail
	if err := r.db.Where("payment_destination_id = ? AND completed_at IS NOT NULL", paymentDestinationId).Find(&pds).Error; err != nil {
		return nil, err
	}
	return pds, nil
}

// FindActiveByPaymentDestinationId 削除されていないPaymentDetailを検索する
func (r PaymentDetailRepository) FindActiveByPaymentDestinationId(paymentDestinationId uint) ([]model.PaymentDetail, error) {
	var pds []model.PaymentDetail
	if err := r.db.Where("payment_destination_id = ? AND deleted_at IS NULL", paymentDestinationId).Find(&pds).Error; err != nil {
		return nil, err
	}
	return pds, nil
}

func (r PaymentDetailRepository) FindByUserId(id uint) ([]model.PaymentDetail, error) {
	var pds []model.PaymentDetail
	if err := r.db.Where("user_id = ?", id).Find(&pds).Error; err != nil {
		return nil, err
	}
	return pds, nil
}

func (r PaymentDetailRepository) Create(pd *model.PaymentDetail) error {
	return r.db.Create(pd).Error
}

func (r PaymentDetailRepository) Update(pd *model.PaymentDetail) error {
	return r.db.Save(pd).Error
}

func (r PaymentDetailRepository) Delete(pd *model.PaymentDetail) error {
	return r.db.Delete(pd).Error
}
