//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"waroka/controllers"
	"waroka/repository"
	"waroka/services"
)

func InitializeUserController(db *gorm.DB) *controllers.UserController {
	wire.Build(
		repository.NewUserRepository,
		services.NewUserService,
		services.NewAuthService,
		controllers.NewUserController)
	return &controllers.UserController{}
}

func InitializeAuthController(db *gorm.DB) *controllers.AuthController {
	wire.Build(
		repository.NewUserRepository,
		services.NewUserService,
		services.NewAuthService,
		controllers.NewAuthController)
	return &controllers.AuthController{}
}

func InitializeRoomController(db *gorm.DB) *controllers.RoomController {
	wire.Build(
		repository.NewRoomRepository,
		repository.NewUserRepository,
		repository.NewPaymentDestinationRepository,
		repository.NewPaymentDetailRepository,
		services.NewRoomService,
		services.NewUserService,
		services.NewPaymentService,
		controllers.NewRoomController)
	return &controllers.RoomController{}
}

func InitializePaymentController(db *gorm.DB) *controllers.RoomController {
	wire.Build(
		repository.NewUserRepository,
		repository.NewPaymentDestinationRepository,
		repository.NewPaymentDetailRepository,
		services.NewRoomService,
		services.NewUserService,
		services.NewPaymentService,
	)
}
