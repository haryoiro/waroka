package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waroka/services"
)

type AuthController struct {
	userService services.IUserService
	authService services.IAuthService
}

func NewAuthController(
	u services.IUserService,
	a services.IAuthService,
) *AuthController {
	return &AuthController{
		userService: u,
		authService: a,
	}
}

func (a *AuthController) RegisterRoutes(e *echo.Echo) {
	auth := e.Group("/auth")

	auth.POST("/signin", a.signin)
	auth.POST("/signup", a.signup)
}

func (a *AuthController) signin(c echo.Context) error {
	var authReq services.AuthRequest
	err := c.Bind(&authReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Idかパスワードが不正です。")
	}

	token, err := a.authService.Signin(authReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Idかパスワードが不正です。")
	}

	cookie := &http.Cookie{
		Name:   "token",
		Value:  *token,
		MaxAge: 604800,
		Path:   "/",
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ログインしました。",
	})
}

func (a *AuthController) signup(c echo.Context) error {
	var authReq services.AuthRequest
	err := c.Bind(&authReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Idかパスワードが不正です。")
	}

	token, err := a.authService.Signup(authReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	cookie := &http.Cookie{
		Name:   "token",
		Value:  *token,
		MaxAge: 604800,
		Path:   "/",
	}

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ログインしました。",
	})
}
