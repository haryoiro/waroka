package repository

import (
	"errors"
	"gorm.io/gorm"
	"waroka/model"
)

type IUserRepository interface {
	Exists(name *string) (bool, error)
	Save(user *model.User) error
	CreateUser(user *model.User) error
	FindAll() ([]model.User, error)
	FindById(id *uint) (*model.User, error)
	FindByName(name *string) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u UserRepository) FindById(id *uint) (*model.User, error) {
	var user model.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) CreateUser(user *model.User) error {
	return u.db.Create(user).Error
}

func (u UserRepository) Exists(name *string) (bool, error) {
	var user model.User
	if err := u.db.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if user.Name == "" {
		return false, nil
	}

	return true, nil
}

func (u UserRepository) Save(user *model.User) error {
	if err := u.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u UserRepository) FindAll() ([]model.User, error) {
	var users []model.User

	//r.db.Table("rooms").
	//	Select("rooms.*").
	//	Joins("left join users_rooms as ur on ur.user_id = rooms.id").
	//	Where("ur.user_id = ?", userId).Find(&rooms).Error

	u.db.Table("users").
		Select("users.*").
		Joins("LEFT JOIN usres_rooms AS ur ON ur.user_id = ")

	if err := u.db.Preload("Rooms").Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return users, nil
}

func (u UserRepository) FindByName(name *string) (*model.User, error) {
	var user model.User
	if err := u.db.Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
