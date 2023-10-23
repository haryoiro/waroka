package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
)

func main() {
    e := echo.New()

    e.Use(middleware.Recover())
    e.Use(middleware.Logger())

    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowCredentials: true,
        AllowOrigins: []string{"http://localhost:8080"},
        AllowMethods: []string{
            http.MethodPost,
            http.MethodGet,
        },
    }))


    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    e.POST("/signup", signup)
    e.POST("/signin", signin)


    e.Logger.Fatal(e.Start(":8080"))
}


type User struct {
    Name string `json:"name" xml:"name" form:"name" query:"name"`
    Password string `json:"password" xml:"password" form:"password" query:"password"`
}


func signup(c echo.Context) error {
    var user User
    err := c.Bind(&user)
    if err != nil {
        return c.String(http.StatusBadRequest, "フォームの値が不正です。")
    }

    // User重複検索

    // password8種か

    // User作成

    cookie := &http.Cookie{
        Name: "user",
        Value: user.Name,
        MaxAge: 604800,
        Path: "/",
    }

    c.SetCookie(cookie)
    return c.String(http.StatusCreated, "登録完了")
}

func signin(c echo.Context) error {
    var user User
    err := c.Bind(&user); if err != nil {
        return c.String(http.StatusBadRequest, "パスワードかIDが不正です")
    }

    // ユーザー検索

    // パスワーdハッシュ化

    // 検証

    cookie := &http.Cookie{
        Name: "user",
        Value: user.Name,
        MaxAge: 608400,
        Path: "/",
    }

    c.SetCookie(cookie)
    return c.String(http.StatusCreated, "ログイン成功")
}


func AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {

        cookie, err := c.Cookie("user")
        if err != nil {
            return c.String(http.StatusNotAcceptable, "認証されていません")
        }

        // ユーザーを検索
        c.Set("user", cookie.Value)

        if err := next(c); err != nil {
            c.Error(err)
        }

        return nil
    }
}