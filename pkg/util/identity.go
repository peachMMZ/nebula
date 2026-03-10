package util

import (
	"math/rand"
)

const (
	// 短ID字符集（纯数字）
	shortIDChars = "0123456789"
	// 短ID长度
	shortIDLength = 6
)

// GenerateShortID 生成短ID
func GenerateShortID() string {
	b := make([]byte, shortIDLength)
	for i := range b {
		b[i] = shortIDChars[rand.Intn(len(shortIDChars))]
	}
	return string(b)
}
