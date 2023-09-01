package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// 加密小写返回
func Md5encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

// 加密大写返回
func Md5StrToUpper(code string) string {
	return strings.ToUpper(Md5encoder(code))
}

// 密码加盐
func SaltPassWord(pw, salt string) string {
	saltPW := fmt.Sprintf("%s$%s", Md5encoder(pw), salt)
	return saltPW
}

// 检验密码
func CheckPassWord(rpw, salt, pw string) bool {
	return pw == SaltPassWord(rpw, salt)
}
