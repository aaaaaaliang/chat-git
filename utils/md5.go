package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(tempStr))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + "------" + password)
	return Md5Encode(plainpwd+salt) == password
}
