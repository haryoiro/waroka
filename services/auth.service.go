package services

import (
	"errors"
	"waroka/model"
	"waroka/repository"
	"waroka/utils"
)

type IAuthService interface {
	Signin(request AuthRequest) (*string, error)
	Signup(request AuthRequest) (*string, error)
}

type AuthService struct {
	userRepo repository.IUserRepository
}

func NewAuthService(repo repository.IUserRepository) IAuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (a AuthService) Signin(request AuthRequest) (*string, error) {
	user, err := a.userRepo.FindByName(&request.Name)
	if err != nil {
		return nil, err
	}

	if !utils.Verify(request.Password, user.Password) {
		return nil, err
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (a AuthService) Signup(request AuthRequest) (*string, error) {
	isExist, err := a.userRepo.Exists(&request.Name)
	if err != nil {
		return nil, errors.New("ユーザーを作成できませんでした。")
	}

	// ユーザーが重複
	if isExist {
		return nil, errors.New("すでにユーザーが存在します。")
	}

	// パスワードをハッシュ化
	hashed := utils.Sign(request.Password)
	if hashed == "" {
		return nil, errors.New("しばらく使用できません。")
	}

	// ユーザーを作成
	user := &model.User{
		Name:     request.Name,
		Password: hashed,
	}

	err = a.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// トークンを作成
	token, err := utils.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

type AuthRequest struct {
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
}
