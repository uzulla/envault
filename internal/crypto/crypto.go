package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	MagicBytes = "ENVAULT1"
	
	KeyLength  = 32 // AES-256-GCM用の32バイト鍵
	NonceSize  = 12 // GCMの標準nonce長
	SaltLength = 16 // 鍵導出用のソルト長
	
	ArgonTime    = 1
	ArgonMemory  = 64 * 1024
	ArgonThreads = 4
)

var (
	ErrInvalidFile      = errors.New("無効なファイル形式です")
	ErrDecryptionFailed = errors.New("復号化に失敗しました。パスワードが間違っている可能性があります")
)

func Encrypt(data []byte, password string) ([]byte, error) {
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("ソルトの生成に失敗しました: %w", err)
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("AESブロック暗号の初期化に失敗しました: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCMモードの初期化に失敗しました: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonceの生成に失敗しました: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	result := []byte(MagicBytes)
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func Decrypt(encryptedData []byte, password string) ([]byte, error) {
	if len(encryptedData) < len(MagicBytes)+SaltLength+NonceSize {
		return nil, ErrInvalidFile
	}

	if string(encryptedData[:len(MagicBytes)]) != MagicBytes {
		return nil, ErrInvalidFile
	}

	offset := len(MagicBytes)
	salt := encryptedData[offset : offset+SaltLength]
	offset += SaltLength

	nonce := encryptedData[offset : offset+NonceSize]
	offset += NonceSize

	ciphertext := encryptedData[offset:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("AESブロック暗号の初期化に失敗しました: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCMモードの初期化に失敗しました: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		ArgonTime,
		ArgonMemory,
		ArgonThreads,
		KeyLength,
	)
}
