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

type UserDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Rooms     []RoomDTO `json:"rooms"`
}
type RoomDTO struct {
	ID       uint   `json:"id"`
	RoomName string `json:"room_name"`
}

func UserToDTO(u model.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func UsersToDTOs(users []model.User) []UserDTO {
	var userDTOs []UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, UserToDTO(user))
	}
	return userDTOs
}

func (u *UserController) getAll(c echo.Context) error {
	users, err := u.userService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "ユーザーを取得できません")
	}

	userDtos := UsersToDTOs(users)

	return c.JSON(http.StatusOK, userDtos)
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

	userDto := UserToDTO(*user)

	return c.JSON(http.StatusOK, userDto)
}

func (u *UserController) create(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := u.userService.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userDto := UserToDTO(*user)

	return c.JSON(http.StatusCreated, userDto)
}
