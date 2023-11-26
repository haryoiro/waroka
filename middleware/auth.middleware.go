package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waroka/utils"
)

func AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("token")
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "ログインして下さい"})
		}

		tokenString := token.Value
		id, err := utils.UserIdFromToken(&tokenString)
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "ログインして下さい"})
		}

		if id == nil {
			return c.JSON(http.StatusNotAcceptable, map[string]interface{}{"message": "ログインして下さい"})
		}

		c.Set("user", *id)

		if err := next(c); err != nil {
			c.Error(err)

		}

		return nil
	}
}
