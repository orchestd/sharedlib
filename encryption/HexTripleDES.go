package encryption

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"unsafe"
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

	err = checkSafetyForCryptBlocks(block, decrypted, data)
	if err != nil {
		return "", err
	}
	decrypter.CryptBlocks(decrypted, data)

	decrypted, err = PKCS5UnPadding(decrypted)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", decrypted), nil
}

// functions copied from alias - internal package golang.org/x/crypto/internal/alias

// AnyOverlap reports whether x and y share memory at any (not necessarily
// corresponding) index. The memory beyond the slice length is ignored.
func aliasAnyOverlap(x, y []byte) bool {
	return len(x) > 0 && len(y) > 0 &&
		uintptr(unsafe.Pointer(&x[0])) <= uintptr(unsafe.Pointer(&y[len(y)-1])) &&
		uintptr(unsafe.Pointer(&y[0])) <= uintptr(unsafe.Pointer(&x[len(x)-1]))
}

// InexactOverlap reports whether x and y share memory at any non-corresponding
// index. The memory beyond the slice length is ignored. Note that x and y can
// have different lengths and still not have any inexact overlap.
//
// InexactOverlap can be used to implement the requirements of the crypto/cipher
// AEAD, Block, BlockMode and Stream interfaces.
func aliasInexactOverlap(x, y []byte) bool {
	if len(x) == 0 || len(y) == 0 || &x[0] == &y[0] {
		return false
	}
	return aliasAnyOverlap(x, y)
}

func checkSafetyForCryptBlocks(b cipher.Block, dst, src []byte) error {
	if len(src)%b.BlockSize() != 0 {
		return fmt.Errorf("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		return fmt.Errorf("crypto/cipher: output smaller than input")
	}
	if aliasInexactOverlap(dst[:len(src)], src) {
		return fmt.Errorf("crypto/cipher: invalid buffer overlap")
	}
	return nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length <= unpadding {
		return nil, fmt.Errorf("key/hash missmatch")
	}
	return origData[:(length - unpadding)], nil
}
