package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type AssetDecryptor struct {
	aesChiper cipher.Block
}

func New() *AssetDecryptor {
	key, _ := base64.StdEncoding.DecodeString("5Mp78iBLX9gVvDjWGCqfbzRHS7hK3JiR")
	newCipher, _ := aes.NewCipher(key)
	return &AssetDecryptor{aesChiper: newCipher}
}

func (desc AssetDecryptor) Encrypt(input []byte) ([]byte, error) {
	blockSize := desc.aesChiper.BlockSize()
	rawData := PKCS7Padding(input, blockSize)
	cipherText := make([]byte, blockSize+len(rawData))
	iv := cipherText[:blockSize]

	mode := cipher.NewCBCEncrypter(desc.aesChiper, iv)
	mode.CryptBlocks(cipherText[blockSize:], input)

	return cipherText, nil
}

func (desc AssetDecryptor) Decrypt(input []byte) ([]byte, error) {
	blockSize := desc.aesChiper.BlockSize()
	iv := input[:blockSize]
	encryptData := input[blockSize:]

	mode := cipher.NewCBCDecrypter(desc.aesChiper, iv)
	mode.CryptBlocks(encryptData, encryptData)
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
