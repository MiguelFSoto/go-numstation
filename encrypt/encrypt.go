package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"bytes"
)

func EncryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad the plaintext to be a multiple of the block size
	plaintext = PKCS7Padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(ciphertext, plaintext)
	hexEncoded := hex.EncodeToString(ciphertext)

	return hexEncoded, nil
}

func DecryptAES(ciphertext string, key []byte) ([]byte, error) {
	decodedCiphertext, _ := hex.DecodeString(ciphertext)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a CBC mode decrypter
	decipher := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))

	// Decrypt the ciphertext
	decrypted := make([]byte, len(decodedCiphertext))
	decipher.CryptBlocks(decrypted, decodedCiphertext)

	// Remove padding
	decrypted = PKCS7Unpadding(decrypted)

	return decrypted, nil
}

// PKCS7Padding pads the given data to be a multiple of blockSize.
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7Unpadding removes PKCS7 padding from the data.
func PKCS7Unpadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}


func test() {
	message := "Hello, AES!"
	key := []byte("0123456789abcdef") // 16-byte key for AES-128

	// Encryption
	ciphertext, err := EncryptAES([]byte(message), key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	// Encode the ciphertext in hexadecimal
	fmt.Printf("Original message: %s\n", message)
	fmt.Printf("Encrypted ciphertext (Hex): %s\n", ciphertext)

	// Decryption
	decrypted, err := DecryptAES(ciphertext, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Printf("Decrypted message: %s\n", decrypted)
}

