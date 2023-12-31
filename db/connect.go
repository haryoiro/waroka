package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func Connect() (*gorm.DB, error) {
	// prduction, local, testで分ける
	// 本番環境では環境変数から取得する

	env := os.Getenv("ENVIRONMENT")

	var dsn string
	if env == "PROD" {
		dsn = "root:wN.=f2m,$,#GfJ[e@tcp(127.0.0.1)/waroka?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true"
	} else {
		dsn = "root:@tcp(127.0.0.1)/waroka?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
