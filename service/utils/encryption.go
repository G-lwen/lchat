package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(values ...string) string {
	h := md5.New()
	for _, value := range values {
		h.Write([]byte(value))
	}
	return hex.EncodeToString(h.Sum(nil))
}
