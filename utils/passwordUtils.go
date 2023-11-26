package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func GetSecret() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("シークレットの読み込みに失敗しました。")
	}

	secret := os.Getenv("HTTP_SECRET")
	return secret
}

func Sign(plain string) string {
	secret := GetSecret()
	fmt.Println("secret " + secret)
	signed, err := getBinaryBySHA256WithKey(plain, secret)
	if err != nil {
		return ""
	}

	return b64.StdEncoding.EncodeToString(signed)
}

func Verify(source string, target string) bool {
	signed := Sign(source)
	return signed == target
}

func getBinaryBySHA256(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

func getBinaryBySHA256WithKey(msg, key string) ([]byte, error) {
	mac := hmac.New(sha256.New, getBinaryBySHA256(key))
	_, err := mac.Write([]byte(msg))
	return mac.Sum(nil), err
}
