package typeConverters

// copy from helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
)

func StringArrayToInterfaceArray(strings []string) []interface{} {
	var intfs []interface{}
	for _, s := range strings {
		intfs = append(intfs, s)
	}
	return intfs
}

func InterfaceToBool(i interface{}) (bool, error) {
	res := false

	switch i.(type) {
	case bool:
		res = i.(bool)
	default:
		return res, errors.New(fmt.Sprint(i) + " is not a boolean.")
	}

	return res, nil
}

func InterfaceToInt64(i interface{}) (int64, error) {
	if res, err := strconv.ParseInt(fmt.Sprint(i), 10, 64); err != nil {
		return res, err
	} else {
		return res, nil
	}
}

func MapInterfaceToMapString(interfaceMap map[string]interface{}) map[string]string {
	stringMap := make(map[string]string)
	for key, val := range interfaceMap {
		stringMap[key] = fmt.Sprint(val)
	}
	return stringMap
}

func EncryptBase64AES(key []byte, text string) (string, error) {
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

func DecryptBase64AES(key []byte, cryptoText string) (string, error) {
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
