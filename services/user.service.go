package services

import (
	"waroka/model"
	"waroka/repository"
)

type IUserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserById(id uint) (*model.User, error)
	CreateUser(user *model.User) error
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserService{
		repo: repo,
	}
}

func (u UserService) GetAllUsers() ([]model.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserService) GetUserById(id uint) (*model.User, error) {
	user, err := u.repo.FindById(&id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserService) CreateUser(user *model.User) error {
	err := u.repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
