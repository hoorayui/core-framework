package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(userID string, password string) string {
	b := md5.Sum([]byte(userID + password))
	return hex.EncodeToString(b[0:])
}
