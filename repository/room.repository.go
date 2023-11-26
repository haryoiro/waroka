package repository

import (
	"errors"
	"gorm.io/gorm"
	"waroka/model"
)

type IRoomRepository interface {
	CreateRoom(room *model.Room) (*model.Room, error)
	JoinRoom(room *model.Room, user *model.User) error
	FindById(id uint) (*model.Room, error)
	FindByUserId(userId uint) ([]model.Room, error)
	TotalUnpaidAmount(roomId uint) (int, error)
}

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r RoomRepository) CreateRoom(room *model.Room) (*model.Room, error) {
	err := r.db.Create(room).Error
	if err != nil {
		return nil, errors.New("部屋の作成に失敗しました")
	}
	return room, nil
}

func (r RoomRepository) JoinRoom(room *model.Room, user *model.User) error {
	if room == nil || user == nil {
		return errors.New("部屋に入室できませんでした")
	}

	var count int64
	r.db.Table("users_rooms").Where("user_id = ? AND room_id = ?", user.ID, room.ID).Count(&count)
	if count > 0 {
		return errors.New("すでに入室しています")
	}

	if err := r.db.Model(&room).Association("Users").Append(&user); err != nil {
		return err
	}

	return nil
}

func (r RoomRepository) FindById(roomId uint) (*model.Room, error) {
	var room model.Room
	if err := r.db.Where("id = ?", roomId).First(&room).Error; err != nil {
		return nil, errors.New("部屋が見つかりませんでした。")
	}

	return &room, nil
}

func (r RoomRepository) FindByUserId(userId uint) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Table("users_rooms").Where("user_id = ?", userId).Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// TotalUnpaidAmount は部屋の未払い金額を返す
func (r RoomRepository) TotalUnpaidAmount(roomId uint) (int, error) {
	var totalAmount int
	err := r.db.Table("payment_destination").
		Select("SUM(payment_detail.amount)").
		Joins("LEFT OUTER JOIN payment_detail ON payment_destination.id = payment_detail.payment_destination_id").
		Where("payment_destination.room_id = ? AND payment_destination.deleted_at IS NULL AND payment_detail.deleted_at IS NULL", roomId).
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}

// TotalAmount は部屋の合計金額を返す
func (r RoomRepository) TotalAmount(roomId uint) (int, error) {
	var totalAmount int

	err := r.db.Table("payment_destination").
		Select("SUM(payment_destination.amount)").
		Where("payment_destination.room_id = ? AND payment_destination.deleted_at IS NULL", roomId).
		Scan(&totalAmount).Error

	if err != nil {
		return 0, err
	}

	return totalAmount, nil
}
