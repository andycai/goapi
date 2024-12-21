package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// MD5Hash 计算字符串的MD5哈希值
func MD5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1Hash 计算字符串的SHA1哈希值
func SHA1Hash(text string) string {
	hash := sha1.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA256Hash 计算字符串的SHA256哈希值
func SHA256Hash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA512Hash 计算字符串的SHA512哈希值
func SHA512Hash(text string) string {
	hash := sha512.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(text string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(text)
}

// GenerateRandomBytes 生成指定长度的随机字节数组
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AESEncrypt AES加密
func AESEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 创建加密器
	stream := cipher.NewCFBEncrypter(block, iv)

	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)

	// 将IV和密文组合
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return result, nil
}

// AESDecrypt AES解密
func AESDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 检查密文长度
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// 提取IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建解密器
	stream := cipher.NewCFBDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// EncryptString 加密字符串
func EncryptString(plaintext string, key []byte) (string, error) {
	ciphertext, err := AESEncrypt([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return Base64Encode(ciphertext), nil
}

// DecryptString 解密字符串
func DecryptString(ciphertext string, key []byte) (string, error) {
	data, err := Base64Decode(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext, err := AESDecrypt(data, key)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GenerateAESKey 生成AES密钥
func GenerateAESKey(bits int) ([]byte, error) {
	if bits != 128 && bits != 192 && bits != 256 {
		return nil, errors.New("invalid key size")
	}
	return GenerateRandomBytes(bits / 8)
}

// CompareHashes 安全比较两个哈希值
func CompareHashes(hash1, hash2 string) bool {
	if len(hash1) != len(hash2) {
		return false
	}

	var result byte
	for i := 0; i < len(hash1); i++ {
		result |= hash1[i] ^ hash2[i]
	}
	return result == 0
}

// HashPassword 对密码进行哈希处理
func HashPassword(password string) string {
	// 生成随机盐值
	salt, err := GenerateRandomBytes(16)
	if err != nil {
		return ""
	}

	// 将密码和盐值组合
	combined := append([]byte(password), salt...)

	// 计算哈希值
	hash := sha256.New()
	hash.Write(combined)
	hashedPassword := hash.Sum(nil)

	// 将盐值和哈希值组合
	result := make([]byte, len(salt)+len(hashedPassword))
	copy(result, salt)
	copy(result[len(salt):], hashedPassword)

	return hex.EncodeToString(result)
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	// 解码哈希值
	decoded, err := hex.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	// 提取盐值和哈希值
	salt := decoded[:16]
	hash := decoded[16:]

	// 计算密码哈希值
	combined := append([]byte(password), salt...)
	h := sha256.New()
	h.Write(combined)
	newHash := h.Sum(nil)

	// 比较哈希值
	return CompareHashes(hex.EncodeToString(hash), hex.EncodeToString(newHash))
}
