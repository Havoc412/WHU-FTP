package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// info 加密算法集合
// EncryptString 使用AES-256加密字符串
func EncryptString_AES256(plaintext, keyString, ivString string) (string, error) {
	// 将字符串表示的密钥和IV转换为字节数组
	key := []byte(keyString)
	iv := []byte(ivString)

	// 新建 cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充原文以符合块大小
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append([]byte(plaintext), bytes.Repeat([]byte{byte(padding)}, padding)...)

	// 加密文本
	ciphertext := make([]byte, len(padtext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padtext)

	// 将密文转换为 Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
