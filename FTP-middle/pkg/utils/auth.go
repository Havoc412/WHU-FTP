package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// generateSessionID creates a new random session ID.
func GenerateSessionID() string {
	// 定义生成的字节长度，这里生成的是 32 字节长度的随机数
	b := make([]byte, 16) // 调整大小以改变会话 ID 的长度
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	// 返回十六进制编码的字符串
	return hex.EncodeToString(b)
}
