package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// MD5 计算字符串的MD5哈希值
func MD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// MD5Bytes 计算字节数组的MD5哈希值
func MD5Bytes(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1 计算字符串的SHA1哈希值
func SHA1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA1Bytes 计算字节数组的SHA1哈希值
func SHA1Bytes(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA256 计算字符串的SHA256哈希值
func SHA256(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA256Bytes 计算字节数组的SHA256哈希值
func SHA256Bytes(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA512 计算字符串的SHA512哈希值
func SHA512(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// SHA512Bytes 计算字节数组的SHA512哈希值
func SHA512Bytes(data []byte) string {
	hash := sha512.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// HMACSHA256 使用指定密钥计算HMAC-SHA256
func HMACSHA256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HMACSHA1 使用指定密钥计算HMAC-SHA1
func HMACSHA1(data, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// EncodeBase64 Base64编码
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 Base64解码
func DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// EncodeURLBase64 URL安全的Base64编码
func EncodeURLBase64(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeURLBase64 URL安全的Base64解码
func DecodeURLBase64(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

// GenerateRandomBytes 生成指定长度的随机字节数组
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString 生成指定长度的随机字符串（十六进制表示）
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length / 2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AESEncrypt AES加密（使用CBC模式和PKCS7填充）
func AESEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// CBC模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedPlaintext := pkcs7Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// 将IV和密文组合
	result := make([]byte, len(iv)+len(ciphertext))
	copy(result, iv)
	copy(result[len(iv):], ciphertext)

	return result, nil
}

// AESDecrypt AES解密（使用CBC模式和PKCS7填充）
func AESDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("密文太短")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 分离IV和密文
	iv := ciphertext[:aes.BlockSize]
	actualCiphertext := ciphertext[aes.BlockSize:]

	// CBC模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(actualCiphertext))
	mode.CryptBlocks(plaintext, actualCiphertext)

	// 去除填充
	return pkcs7Unpadding(plaintext)
}

// pkcs7Padding 对数据进行PKCS7填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpadding 移除PKCS7填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据长度为0")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("无效的填充")
	}

	// 验证填充是否正确
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("无效的填充")
		}
	}

	return data[:length-padding], nil
}

// AESGCMEncrypt AES-GCM加密（AEAD认证加密）
func AESGCMEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 生成随机NonceSize
	nonce := make([]byte, 12) // GCM推荐的Nonce大小
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// GCM模式
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 加密并认证
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// 将Nonce和密文组合
	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result, nonce)
	copy(result[len(nonce):], ciphertext)

	return result, nil
}

// AESGCMDecrypt AES-GCM解密（AEAD认证解密）
func AESGCMDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("密文太短")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// GCM模式
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 分离Nonce和密文
	nonce := ciphertext[:12]
	actualCiphertext := ciphertext[12:]

	// 解密并验证
	plaintext, err := aesgcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败：%v", err)
	}

	return plaintext, nil
}
