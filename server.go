package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	d "waroka/db"
	"waroka/di"
	"waroka/model"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	db, err := d.Connect()
	if err != nil {
		panic("MySQLへの接続に失敗しました。")
	}

	db.Exec("SHOW ENGINE INNODB STATUS")

	err = db.AutoMigrate(&model.User{}, &model.Room{}, &model.PaymentDestination{}, &model.PaymentDetail{})
	if err != nil {
		panic("マイグレーションに失敗しました。")
	}

	db.Exec("SET FOREIGN_KEY_CHECKS=1")

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//   AllowCredentials: true,
	//   AllowOrigins: []string{"http://localhost:8080"},
	//   AllowMethods: []string{
	//       http.MethodPost,
	//       http.MethodGet,
	//   },
	//}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	authController := di.InitializeAuthController(db)
	authController.RegisterRoutes(e)

	userController := di.InitializeUserController(db)
	userController.RegisterRoutes(e)

	roomController := di.InitializeRoomController(db)
	roomController.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
	defer e.Logger.Fatal(e.Close())
}
