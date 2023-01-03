package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"unicode"
)

// 简单计算一个 user password 作为token
func HmacAccessToken(user string, password string) []byte {

	hnew := hmac.New(sha256.New224, []byte(fmt.Sprintf("%v%v", user, password)))
	return hnew.Sum(nil)
}

// 参数字符检查
func VerifyLetter(param string) bool {

	for _, v := range param {
		if !unicode.IsLetter(v) {
			return false
		}
	}
	return true
}
