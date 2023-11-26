package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
	"waroka/model"
)

var secretKey = os.Getenv("HTTP_SECRET")

type jwtClaims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}

func CreateToken(user *model.User) (string, error) {
	claims := &jwtClaims{
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func UserIdFromToken(tokenString *string) (*uint, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("トークンが不正です。")
	}

	mapClaims := token.Claims.(jwt.MapClaims)
	if mapClaims["id"] != nil {
		floatId := mapClaims["id"].(float64)
		uintId := uint(floatId)
		fmt.Println("user id ", uintId)
		return &uintId, nil
	} else {
		return nil, errors.New("トークンが不正です。")
	}
}
