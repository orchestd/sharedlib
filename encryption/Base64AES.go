package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(data, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("the key cannot be empty")
	} else if len(key) > 16 {
		return "", fmt.Errorf("the key must be smaller than 16 characters %v", key)
	} else {
		var res string
		if tmp, err := encryptBase64AES([]byte(key), data); err != nil {
			return "", fmt.Errorf("cannot encrypt data. " + err.Error())
		} else {
			res = tmp
		}
		return res, nil
	}
}

func encryptBase64AES(key []byte, text string) (string, error) {
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("the key cannot be empty")
	} else if len(key) > 16 {
		return "", fmt.Errorf("the key must be smaller than 16 characters %v", key)
	} else {
		var res string
		if tmp, err := decryptBase64AES([]byte(key), encrypted); err != nil {
			return "", fmt.Errorf("cannot decrypt data. " + err.Error())
		} else {
			res = tmp
		}
		return res, nil
	}
}

func decryptBase64AES(key []byte, cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
