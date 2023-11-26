package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"waroka/middleware"
	"waroka/model"
	"waroka/services"
)

type UserController struct {
	userService services.IUserService
	authService services.IAuthService
}

func NewUserController(
	u services.IUserService,
	a services.IAuthService,
) *UserController {
	return &UserController{
		userService: u,
		authService: a,
	}
}

func (u *UserController) RegisterRoutes(e *echo.Echo) {

	user := e.Group("/users")

	user.Use(middleware.AuthenticateMiddleware)

	user.GET("", u.getAll)
	user.POST("", u.create)
	user.GET("/:id", u.getById)
}

func (u *UserController) getAll(c echo.Context) error {
	users, err := u.userService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ユーザーを取得できません")
	}
	return c.JSON(http.StatusOK, users)
}

func (u *UserController) getById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Id")
	}

	user, err := u.userService.GetUserById(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ユーザーが取得できませんでした")
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, "ユーザーが見つかりませんでした")
	}
	return c.JSON(http.StatusOK, user)
}

func (u *UserController) create(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := u.userService.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}
