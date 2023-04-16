package encryption

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
)

func TripleDesEncrypt(keyStr, dataStr string) (string, error) {
	key := []byte(keyStr)
	data := []byte(dataStr)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	iv := key[:des.BlockSize]
	origData := PKCS5Padding(data, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	mode.CryptBlocks(encrypted, origData)
	return fmt.Sprintf("%x", encrypted), nil
}

func TripleDesDecrypt(keyStr, dataStr string) (string, error) {
	key := []byte(keyStr)
	data, err := hex.DecodeString(dataStr)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}

	iv := key[:des.BlockSize]
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))

	decrypted, err = decrypt(decrypter, decrypted, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", decrypted), nil
}

func decrypt(decrypter cipher.BlockMode, decrypted, data []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovered from panic, library: sharedlib, func: decrypt, %v", r)
		}
	}()

	decrypter.CryptBlocks(decrypted, data)
	decrypted = PKCS5UnPadding(decrypted)

	return decrypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
